package handler

import (
	"github.com/jhamiltonjunior/rinha-de-backend/app/worker"
	"github.com/valyala/fasthttp"
)

func Payments(ctx *fasthttp.RequestCtx) {
	bodyCopy := make([]byte, len(ctx.PostBody()))
	copy(bodyCopy, ctx.PostBody())

	paymentWorker := worker.PaymentWorker{
		Body: bodyCopy,
	}

	select {
	case worker.SegureOChann <- paymentWorker:
		sendJSONResponse(ctx, fasthttp.StatusAccepted)
	default:
		sendJSONResponse(ctx, fasthttp.StatusTooManyRequests)
	}
}
