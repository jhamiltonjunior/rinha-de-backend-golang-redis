package handler

import (
	"github.com/jhamiltonjunior/rinha-de-backend/app/worker"
	"github.com/valyala/fasthttp"
)

func Payments(ctx *fasthttp.RequestCtx) {
	sendJSONResponse(ctx, fasthttp.StatusAccepted)
	go worker.ProcessPayment(ctx.PostBody())
}
