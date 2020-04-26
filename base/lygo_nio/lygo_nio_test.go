package lygo_nio

import (
	"fmt"
	"testing"
	"time"
)

func TestSimple(t *testing.T) {

	server := NewNioServer(10001)
	err := server.Open()
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	server.OnMessage(onMessage)

	client := NewNioClient("localhost", 10001)
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
	fmt.Println("Response from server", response)

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

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func onMessage(message *NioMessage) interface{}{
	fmt.Println("MESSAGE GOT FROM SERVER", message)
	return "custom response from server handled message"
}
