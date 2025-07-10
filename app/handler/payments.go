package handler

import (
	"encoding/json"
	"fmt"

	"github.com/jhamiltonjunior/rinha-de-backend/app/worker"
	"github.com/valyala/fasthttp"
)

func Payments(ctx *fasthttp.RequestCtx) {
	
	payment := getJSONBody(ctx)
	paymentBytes, err := json.Marshal(payment)
	if err != nil {
		sendJSONResponse(ctx, fasthttp.StatusBadRequest)
		fmt.Fprintf(ctx, `{"message": "Erro ao serializar o pagamento"}`)
		return
	}

	fmt.Println("Recebendo pagamento:", string(paymentBytes))

	go worker.ProcessPayment(paymentBytes)

	sendJSONResponse(ctx, fasthttp.StatusAccepted)

	fmt.Fprintf(ctx, `{"message": "Bem-vindo à Rinha de Backend!", "data": %s}`)
}