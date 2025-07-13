package handler

import (
	"fmt"

	"github.com/jhamiltonjunior/rinha-de-backend/app/database"
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

func PaymentsSummary(ctx *fasthttp.RequestCtx) {
	from := ctx.QueryArgs().Peek("from")
	to := ctx.QueryArgs().Peek("to")

	if len(from) == 0 {
		from = []byte("1970-01-01T00:00:00Z")
	}

	if len(to) == 0 {
		to = []byte("9999-12-31T23:59:59Z")
	}

	payments, err := database.GetPaymentHistory(database.MongoClient, string(from), string(to))
	if err != nil {
		sendJSONResponse(ctx, fasthttp.StatusInternalServerError)
		return
	}
	
	fmt.Fprintf(ctx, "{\"total\": %d\n, \"payments\": %v}", len(payments), payments)

	sendJSONResponse(ctx, fasthttp.StatusOK)
}