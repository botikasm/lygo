package lygo_http_client

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/valyala/fasthttp"
	"log"
	"testing"
)

func TestSimple(t *testing.T) {

	client := new(HttpClient)

	// Fetch google page via local proxy.
	fmt.Println("https://google.com/")
	code, body, err := client.Get("https://google.com/")
	if code != fasthttp.StatusOK {
		log.Fatalf("Unexpected status code: %d. Expecting %d", code, fasthttp.StatusOK)
	}
	if err != nil {
		log.Fatalf("Error when loading google page through local proxy: %s", err)
	}
	useResponseBody(body)

	// Fetch foobar page via local proxy. Reuse body buffer.
	fmt.Println("https://botika.ai/")
	code, body, err = client.Get("https://botika.ai/")
	if code != fasthttp.StatusOK {
		log.Fatalf("Unexpected status code: %d. Expecting %d", code, fasthttp.StatusOK)
	}
	if err != nil {
		log.Fatalf("Error when loading google page through local proxy: %s", err)
	}
	useResponseBody(body)

}

func TestDownload(t *testing.T) {
	client := new(HttpClient)

	file := "https://gianangelogeminiani.me/download/architecture.png"
	fmt.Println(file)
	code, body, err := client.Get(file)
	if code != fasthttp.StatusOK {
		log.Fatalf("Unexpected status code: %d. Expecting %d", code, fasthttp.StatusOK)
	}
	if err != nil {
		log.Fatalf("Error when loading google page through local proxy: %s", err)
	}

	fileName := "./architecture.png"
	lygo_io.WriteBytesToFile(body, fileName)

	fmt.Println(fileName)
}
//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func useResponseBody(body []byte) {
	fmt.Println(string(body))
}
