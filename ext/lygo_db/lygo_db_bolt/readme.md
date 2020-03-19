#lygo_db_bbolt
##BBolt embedded database wrapper

###Dependencies
This module depend on [bbolt](https://github.com/etcd-io/bbolt) 


`go get go.etcd.io/bbolt/...`

Bolt is a pure Go key/value store inspired by [Howard Chu's][hyc_symas]
[LMDB project][lmdb]. The goal of the project is to provide a simple,
fast, and reliable database for projects that don't require a full database
server such as Postgres or MySQL.

Since Bolt is meant to be used as such a low-level piece of functionality,
simplicity is key. The API will be small and only focus on getting values
and setting values. That's it.


[hyc_symas]: https://twitter.com/hyc_symas
[lmdb]: http://symas.com/mdb/

