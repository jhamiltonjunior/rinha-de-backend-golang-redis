package handler

import (
	"github.com/valyala/fasthttp"
)

func sendJSONResponse(ctx *fasthttp.RequestCtx, statusCode int) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)
}
