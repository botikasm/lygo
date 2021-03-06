package lygo_nio

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_events"
	"strings"
	"testing"
	"time"
)

func TestRunServer(t *testing.T) {
	server := NewNioServer(10001)
	err := server.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	server.OnMessage(onMessage)

	// stop after 3 seconds
	go func() {
		time.Sleep(5 * time.Second)
		server.Close()
	}()

	fmt.Println("Http listening on port:", server.port)
	server.Join()
	fmt.Println("Http CLOSE")
}

func TestRunClientPing(t *testing.T) {
	server := NewNioServer(10001)
	err := server.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("Http listening on port:", server.port)

	// start and stop server every 3 secs
	go func() {
		for {
			time.Sleep(5 * time.Second)
			if server.IsOpen() {
				fmt.Println("SERVER", "off")
				server.Close()
			} else {
				fmt.Println("SERVER", "on")
				server.Open()
			}
		}
	}()

	client := NewNioClient("localhost", 10001)
	client.Secure = true // enable cryptography
	client.OnConnect(func(e *lygo_events.Event) {
		fmt.Println("CLIENT", "connected")
	})
	client.OnDisconnect(func(e *lygo_events.Event) {
		fmt.Println("CLIENT", "disconnected")
	})
	err = client.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	client.Join()
}

func TestSimple(t *testing.T) {

	server := NewNioServer(10001)
	err := server.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	server.OnMessage(onMessage)

	client := NewNioClient("localhost", 10001)
	client.Secure = true // enable cryptography
	err = client.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	client2 := NewNioClient("localhost", 10001)
	err = client2.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	time.Sleep(1 * time.Second)
	fmt.Println("Clients", server.ClientsCount(), server.ClientsId())
	// disconnect second client
	err = client2.Close()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	time.Sleep(1 * time.Second)
	fmt.Println("Clients", server.ClientsCount())

	//-- MESSAGE --//
	response, err := client.Send("hello")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	body := lygo_conv.ToString(response.Body)
	fmt.Println("Response from server:")
	fmt.Println(body)

	// disconnect client
	err = client.Close()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println("exiting....")
	time.Sleep(3 * time.Second)
	fmt.Println("Clients", server.ClientsCount())
}

func TestComplexData(t *testing.T) {
	server := NewNioServer(10001)
	err := server.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	server.OnMessage(onMessage)

	client := NewNioClient("localhost", 10001)
	client.Secure = true // enable cryptography
	err = client.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	data := &map[string]interface{}{
		"name":    "Mario",
		"surname": "Rossi",
		"age":     53,
		"phone":   "+3912214235356",
		"sons":    []string{"Maria", "John"},
	}

	//-- MESSAGE --//
	response, err := client.Send(data)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	body := lygo_conv.ToString(response.Body)
	fmt.Println("Response from server:")
	fmt.Println(body)

	// disconnect client
	err = client.Close()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println("exiting....")
}

func TestBigData(t *testing.T) {
	server := NewNioServer(10001)
	err := server.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	server.OnMessage(onMessage)

	timeStart := time.Now()

	client := NewNioClient("localhost", 10001)
	client.Secure = false // enable cryptography
	err = client.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	bigArray := make([]int, 100000)
	for i := 0; i < len(bigArray); i++ {
		bigArray[i] = i
	}
	data := &map[string]interface{}{
		"name":    "Mario",
		"surname": "Rossi",
		"age":     53,
		"phone":   "+3912214235356",
		"sons":    []string{"Maria", "John"},
		"big":     bigArray,
	}

	//-- MESSAGE --//
	response, err := client.Send(data)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	body := lygo_conv.ToString(response.Body)
	fmt.Println("Response from server:")
	fmt.Println(body)
	response, err = client.Send(data)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	body = lygo_conv.ToString(response.Body)
	fmt.Println("Response from server:")
	fmt.Println(body)

	// disconnect client
	err = client.Close()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	timeEnd := time.Now()
	fmt.Println("elapsed", timeEnd.Sub(timeStart))

	fmt.Println("exiting....")
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func onMessage(message *NioMessage) interface{} {
	body := lygo_conv.ToString(message.Body)
	if strings.Index(body, "{") > -1 {
		m := lygo_conv.ToMap(body)
		if v, b := m["big"]; b {
			a, _ := v.([]interface{})
			fmt.Println("COMPLEX MESSAGE ARRIVED ON SERVER. big:", len(a))
			fmt.Println("Sending response big... ")
			return a
		} else {
			fmt.Println("COMPLEX MESSAGE ARRIVED ON SERVER:", m)
			fmt.Println("Sending response as custom map... ")
			return &map[string]interface{}{
				"tag":    "COMPLEX RESPONSE",
				"body":   body,
				"object": m,
			}
		}
	} else {
		fmt.Println("MESSAGE ARRIVED ON SERVER:", body)
		return "custom response from server handled message"
	}
}
