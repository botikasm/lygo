package lygo_http_client

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/valyala/fasthttp"
	"log"
	"net/url"
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

func TestTinyURL(t *testing.T) {
	urlFull := "http://localhost:63343/ritiro_io_client/index.html?_ijt=qbk2r2ocijg43343og9ivnvr4o#!/02_viewer/menu/ee63b7f4-1766-487e-8762-3a2710320158/04eff121-d533-edc9-7fc2-ebc393895250"
	client := new(HttpClient)
	callUrl := "http://tinyurl.com/api-create.php?url=" + url.QueryEscape(urlFull)
	_, response, err := client.Get(callUrl)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(string(response))
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
