package lygo_http_server

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"expvar"
	"fmt"
	"github.com/botikasm/lygo/base/lygo_strings"
	"github.com/valyala/fasthttp"
	"math/big"
	"sync"
	"time"
)

var (
	errorInvalidConfiguration = errors.New("Configuration is not valid")
)

const maxErrors = 100

// Various counters - see https://golang.org/pkg/expvar/ for details.
var (
	// Counter for total number of fs calls
	fsCalls = expvar.NewInt("fsCalls")

	// Counters for various response status codes
	fsOKResponses          = expvar.NewInt("fsOKResponses")
	fsNotModifiedResponses = expvar.NewInt("fsNotModifiedResponses")
	fsNotFoundResponses    = expvar.NewInt("fsNotFoundResponses")
	fsOtherResponses       = expvar.NewInt("fsOtherResponses")

	// Total size in bytes for OK response bodies served.
	fsResponseBodyBytes = expvar.NewInt("fsResponseBodyBytes")
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type HttpServer struct {
	Config *HttpServerConfig

	//-- hooks --//
	CallbackError ErrorHookCallback

	//-- private --//
	started   bool
	stopped   bool
	fsHandler fasthttp.RequestHandler
	errors    []error
	muxError  sync.Mutex
}

type HttpServerError struct {
	Server  *HttpServer
	Message string
	Error   error
	Context *fasthttp.RequestCtx
}

type ErrorHookCallback func(*HttpServerError)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewHttpServer(config *HttpServerConfig) *HttpServer {
	instance := new(HttpServer)
	instance.Config = config
	instance.stopped = false
	instance.started = false

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpServer) Start() error {
	if !instance.started {
		if nil == instance.Config {
			return errorInvalidConfiguration
		}
		instance.started = true
		instance.stopped = false

		instance.initConfig()
		instance.initFileServer()
		instance.initFileStat()
		instance.serve()

	}
	return nil
}

func (instance *HttpServer) Join() {
	if !instance.stopped {
		for !instance.stopped {
			time.Sleep(10 * time.Second)
		}
	}
}

func (instance *HttpServer) Stop() {
	if !instance.stopped {
		instance.stopped = true
		instance.started = false

	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpServer) initConfig() {
	config := instance.Config
	if nil != config {
		if nil == config.IndexNames || len(config.IndexNames) == 0 {
			config.IndexNames = append(config.IndexNames, "index.html")
		}
	}
}

func (instance *HttpServer) initFileServer() {
	config := instance.Config
	if len(config.FileServerRoot) > 0 && config.FileServerEnabled {
		fs := &fasthttp.FS{
			Root:               config.FileServerRoot,
			GenerateIndexPages: false,
			Compress:           false,
			AcceptByteRange:    false,
		}
		fs.IndexNames = append(fs.IndexNames, config.IndexNames...)
		if config.VHost {
			fs.PathRewrite = fasthttp.NewVHostPathRewriter(0)
		}
		instance.fsHandler = fs.NewRequestHandler()
	}
}

func (instance *HttpServer) initFileStat() {
	if instance.Config.StatEnabled {

	}
}

func (instance *HttpServer) serve() {
	config := instance.Config
	if len(config.Address) > 0 {
		// http
		go instance.listenAndServe(config.Address)
	}
	if len(config.AddressTLS) > 0 {
		// https
		go instance.listenAndServeTLS(config.AddressTLS, config.SslCert, config.SslKey)
	}
}

func (instance *HttpServer) listenAndServe(addr string) {
	if err := fasthttp.ListenAndServe(addr, instance.onHandle); err != nil {
		// log.Fatalf("error in ListenAndServe: %s", err)
		instance.doError(lygo_strings.Format("Error Opening channel: '%s'", instance.Config.Address), err, nil)
	}
}

func (instance *HttpServer) listenAndServeTLS(addr string, certFile string, keyFile string) {
	if len(certFile) > 0 && len(keyFile) > 0 {
		// use files
		if err := fasthttp.ListenAndServeTLS(addr, certFile, keyFile, instance.onHandle); err != nil {
			instance.doError(lygo_strings.Format("Error Opening TLS channel: '%s'", instance.Config.AddressTLS), err, nil)
		}
	} else {
		// auto-generate cert
		cert, key, err := GenerateCert(addr)
		if nil == err {
			if err = fasthttp.ListenAndServeTLSEmbed(addr, cert, key, instance.onHandle); err != nil {
				instance.doError("Error Opening TLS channel", err, nil)
			}
		} else {
			instance.doError("Error Generating Certificate", err, nil)
		}
	}
}

func (instance *HttpServer) doError(message string, err error, ctx *fasthttp.RequestCtx) {
	instance.muxError.Lock()
	go func() {
		defer instance.muxError.Unlock()
		if nil != instance.errors && len(instance.errors) > maxErrors {
			// reset errors
			instance.errors = make([]error, 0)
		}
		instance.errors = append(instance.errors, err)
		if nil != instance.CallbackError {
			instance.CallbackError(&HttpServerError{
				Server:  instance,
				Message: message,
				Context: ctx,
				Error:   err,
			})
		}
	}()
}

func (instance *HttpServer) onHandle(ctx *fasthttp.RequestCtx) {
	instance.stat(ctx)
	path := ctx.Path()
	if len(path) > 0 {

		if instance.Config.StatEnabled && string(path) == "/stats" {
			ExpvarHandler(ctx)
		} else {

			fmt.Println(string(ctx.Host()))

			// file server
			if nil != instance.fsHandler {
				instance.fsHandler(ctx)
			}
		}
	}
}

func (instance *HttpServer) stat(ctx *fasthttp.RequestCtx) {
	if instance.Config.StatEnabled {
		// Increment the number of fsHandler calls.
		fsCalls.Add(1)

		status := ctx.Response.StatusCode()
		switch status {
		case fasthttp.StatusOK:
			fsOKResponses.Add(1)
			fsResponseBodyBytes.Add(int64(ctx.Response.Header.ContentLength()))
		case fasthttp.StatusNotModified:
			fsNotModifiedResponses.Add(1)
		case fasthttp.StatusNotFound:
			fsNotFoundResponses.Add(1)
		default:
			fsOtherResponses.Add(1)
		}
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	s t a t i c
//----------------------------------------------------------------------------------------------------------------------

// GenerateCert generates certificate and private key based on the given host.
func GenerateCert(host string) ([]byte, []byte, error) {
	priv, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, nil, err
	}

	cert := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"I have your data"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().Add(365 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth, x509.ExtKeyUsageClientAuth},
		SignatureAlgorithm:    x509.SHA256WithRSA,
		DNSNames:              []string{host},
		BasicConstraintsValid: true,
		IsCA:                  true,
	}

	certBytes, err := x509.CreateCertificate(
		rand.Reader, cert, cert, &priv.PublicKey, priv,
	)

	p := pem.EncodeToMemory(
		&pem.Block{
			Type:  "PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(priv),
		},
	)

	b := pem.EncodeToMemory(
		&pem.Block{
			Type:  "CERTIFICATE",
			Bytes: certBytes,
		},
	)

	return b, p, err
}
