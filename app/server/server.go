package server

import (
	"fmt"

	"github.com/jhamiltonjunior/rinha-de-backend/app/handler"
	"github.com/valyala/fasthttp"
)

func ListenAndServe(appPort string) {
	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/payments":
			if ctx.IsPost() {
				handler.Payments(ctx)
				return
			}
		case "/payments-summary":
			if ctx.IsGet() {
				handler.PaymentsSummary(ctx)
				return
			}
		case "/purge-payments":
			if ctx.IsPost() {
				handler.PaymentsPurge(ctx)
				return
			}
		}

		ctx.Error("Not Found", fasthttp.StatusNotFound)
	}

	fmt.Println("Servidor rodando em http://localhost:" + appPort)

	if err := fasthttp.ListenAndServe(":"+appPort, requestHandler); err != nil {
		fmt.Printf("Error in ListenAndServe: %s\n", err)
	}
}
