package middlewares

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method

		ctx.Header("Access-Control-Allow-Origin", "*")
		ctx.Header(
			"Access-Control-Allow-Headers",
			"Content-Type, AccessToken, X-CSRF-Token, Authorization, Token, x-token",
		)
		ctx.Header(
			"Access-Control-Allow-Methods",
			"POST, GET, DELETE, PATCH, PUT",
		)
		ctx.Header(
			"Access-Control-Expose-Headers",
			"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type",
		)
		ctx.Header("Access-Control-Allow-Credentials", "true")

		if method == "OPTION" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
	}
}
