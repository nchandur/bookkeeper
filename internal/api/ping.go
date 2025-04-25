package api

import (
	"bookkeeper/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
)

func PingHandler(r *gin.Engine) {
	r.GET("/ping", func(ctx *gin.Context) {
		if err := db.Client.Ping(ctx, nil); err != nil {
			ctx.JSON(http.StatusServiceUnavailable, gin.H{
				"error": err.Error(),
			})
		}

		ctx.JSON(http.StatusOK, gin.H{
			"body":  "pong",
			"error": nil,
		})

	})
}
