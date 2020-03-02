#LyGo

LyGo is a simple application library written in Go.

#Modules

LyGo uses [Modules](https://blog.golang.org/using-go-modules) for dependencies.

To list current module and all his dependency use:
`go list -m all`

To add LyGo as a dependency use:
`go get github.com/botikasm/lygo`

To remove unused dependency use:
`go mod tidy`
