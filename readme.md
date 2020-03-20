LyGo 0.9.16
-

LyGo is a simple application library written in Go.

Current version is a development version in alpha 0.x.

Modules
-

LyGo uses [Modules](https://blog.golang.org/using-go-modules) for dependencies.

To list current module and all his dependency use:
`go list -m all`

To add LyGo as a dependency use:

`go get github.com/botikasm/lygo`

`go get github.com/botikasm/lygo@v0.9.16`

To remove unused dependency use:

`go mod tidy`

To download required dependencies use:

`go build` or `go test`

Version Tagging
-
To tag a version use:

`git tag v0.9.16` 

`git push origin v0.9.16`