#lygo NIO
lygo_nio is a simple network library for client/server 
communication. 
It uses [GOB](https://golang.org/pkg/encoding/gob/) for messages serialisation.

## Cryptography
lygo_nio uses hybrid session encryption with PrivateKey and PublicKey.

## Handshake
If Security is enabled, during handshake client and server share their 
PublicKey used to encrypt the Session Key.

SessionKey are used to encrypt and decrypt message body.

## Transmission Protocol
It uses [GOB](https://golang.org/pkg/encoding/gob/).

## Sample Usage

This example creates a server and a client.
The client start pinging the server and handle "connect" and "disconnect"
events.

```
    server := NewNioServer(10001)
	err := server.Open()
	if nil != err {
		panic(err)
	}
	fmt.Println("Server listening on port:", server.port)

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
		panic(err)
	}
	client.Join()
```
