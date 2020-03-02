LyGo 0.9.4
-

LyGo is a simple application library written in Go.

Modules
-

LyGo uses [Modules](https://blog.golang.org/using-go-modules) for dependencies.

To list current module and all his dependency use:
`go list -m all`

To add LyGo as a dependency use:

`go get github.com/botikasm/lygo`

`go get github.com/botikasm/lygo@v0.9.4`

To remove unused dependency use:
`go mod tidy`

Version Tagging
-
To tag a version use:

`git tag v0.9.4` 

`git push origin v0.9.4`