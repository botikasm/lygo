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
