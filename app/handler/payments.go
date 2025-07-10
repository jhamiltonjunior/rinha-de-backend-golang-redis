package handler

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/jhamiltonjunior/rinha-de-backend/app/utils"
	"github.com/valyala/fasthttp"
)

func Payments(ctx *fasthttp.RequestCtx) {
	
	payment := getJSONBody(ctx)
	payment["requestedAt"] = time.Now().UTC()

	paymentBytes, err := json.Marshal(payment)
	if err != nil {
		sendJSONResponse(ctx, fasthttp.StatusBadRequest)
		fmt.Fprintf(ctx, `{"message": "Erro ao serializar o pagamento"}`)
		return
	}

	fmt.Println("Recebendo pagamento:", string(paymentBytes))

	rever o request talvez seja melhor usar o http.Post do net/http

	body, statusCode := utils.Request("POST", paymentBytes)
	if statusCode < 200 || statusCode > 299 {
		sendJSONResponse(ctx, statusCode)
		fmt.Fprintf(ctx, `{"message": "Erro ao processar o pagamento", "data": %s}`, body)
		return
	}

	sendJSONResponse(ctx, statusCode)

	fmt.Fprintf(ctx, `{"message": "Bem-vindo à Rinha de Backend!", "data": %s}`, body)
}