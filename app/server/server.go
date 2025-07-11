package server

import (
	"fmt"

	"github.com/fasthttp/router"
	"github.com/jhamiltonjunior/rinha-de-backend/app/handler"
	"github.com/valyala/fasthttp"
)
func ListenAndServe(appPort, natsURL string) {
	r := router.New()

	r.POST("/payments", handler.Payments)
	// r.GET("/payments-summary", handler.Payments)
	// r.POST("/purge-payments", handler.Payments)

	fmt.Println("Servidor rodando em http://localhost:" + appPort)
	fasthttp.ListenAndServe(":" + appPort, r.Handler)
}