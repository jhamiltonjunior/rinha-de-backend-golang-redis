package handler

import (
	"encoding/json"
	"fmt"

	"github.com/valyala/fasthttp"
)

func sendJSONResponse(ctx *fasthttp.RequestCtx, statusCode int) {
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(statusCode)
}

func getJSONBody(ctx *fasthttp.RequestCtx) map[string]interface{} {
	body := ctx.PostBody()

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		ctx.SetStatusCode(fasthttp.StatusBadRequest)
		fmt.Fprintf(ctx, "Erro ao decodificar JSON: %v", err)
		return nil
	}
	return data
}
