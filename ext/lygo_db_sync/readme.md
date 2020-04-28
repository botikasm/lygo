# DB Sync 
## Database Synchronization

lygo_db_sync is a client/server TCP/IP database synchronizer from local DB to remote DB.

## Dependencies

### Scripting
To maximize freedom in synchronization behaviour, DB DYNC uses javascript
engine and avoid complex declarative XML like or JSON structures.

This module depend on [Goja](https://github.com/dop251/goja) 

`go get github.com/dop251/goja`

### Arango DB

This module depend on [go-driver](https://github.com/arangodb/go-driver) 

`go get github.com/arangodb/go-driver`
 
 `go get github.com/arangodb/go-driver/http`
 



