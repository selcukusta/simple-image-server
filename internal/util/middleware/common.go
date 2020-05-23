package middleware

import "github.com/valyala/fasthttp"

//CommonMiddleware is using to add common headers to the response
func CommonMiddleware(h fasthttp.RequestHandler) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Add("Access-Control-Allow-Origin", "*")
		ctx.Response.Header.Add("Access-Control-Allow-Methods", "GET")
		h(ctx)
	}
}
