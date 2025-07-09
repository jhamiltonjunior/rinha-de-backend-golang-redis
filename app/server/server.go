package server

import (
	"fmt"

	"github.com/fasthttp/router"
	"github.com/jhamiltonjunior/rinha-de-backend/app/handler"
	"github.com/valyala/fasthttp"
)
func ListenAndServe(appPort, postgresURL, natsURL string) {
	r := router.New()

	r.POST("/payments", handler.HomeHandler)
	r.GET("/payments-summary", handler.HomeHandler)
	r.POST("/purge-payments", handler.HomeHandler)

	fmt.Println("Servidor rodando em http://localhost:" + appPort)
	fasthttp.ListenAndServe(":" + appPort, r.Handler)
}