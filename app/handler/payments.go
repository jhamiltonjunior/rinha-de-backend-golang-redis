package handler

import (
	"fmt"

	"github.com/valyala/fasthttp"
)

func HomeHandler(ctx *fasthttp.RequestCtx) {
	sendJSONResponse(ctx, fasthttp.StatusOK)
	fmt.Fprintf(ctx, `{"message": "Bem-vindo à Rinha de Backend!"}`)
}