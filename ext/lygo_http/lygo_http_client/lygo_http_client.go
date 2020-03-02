package lygo_http_client

import (
	"github.com/valyala/fasthttp"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const methodGet = "GET"
const methodPost = "POST"
const methodPut = "PUT"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type HttpClient struct {
	client *fasthttp.Client
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewHttpClient() *HttpClient {
	instance := new(HttpClient)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpClient) Get(url string) (int, []byte, error) {
	return instance.GetTimeout(url, time.Second*15)
}

func (instance *HttpClient) GetTimeout(url string, timeout time.Duration) (int, []byte, error) {
	return instance.get(url, timeout)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *HttpClient) init() *fasthttp.Client {
	if nil == instance.client {
		instance.client = new(fasthttp.Client)
		instance.client.Name = "lygo_http_client"
	}
	return instance.client
}

func (instance *HttpClient) do(method string, uri string, timeout time.Duration) (statusCode int, body []byte, err error) {

	// request
	req := fasthttp.AcquireRequest()
	req.Header.SetMethod(method)
	req.SetRequestURI(uri)

	// response
	res := fasthttp.AcquireResponse()

	defer func() {
		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(res)
	}()

	client := instance.init()
	err = client.DoTimeout(req, res, timeout)

	return res.StatusCode(), res.Body(), err
}

func (instance *HttpClient) get(uri string, timeout time.Duration) (statusCode int, body []byte, err error) {
	client := instance.init()
	return client.GetTimeout(nil, uri, timeout)
}
