package lygo_http_server_types

import (
	"github.com/gofiber/fiber"
)

type HttpServerError struct {
	Sender  interface{}
	Message string
	Error   error
	Context *fiber.Ctx
}

type CallbackError func(*HttpServerError)
type CallbackLimitReached func(ctx *fiber.Ctx)