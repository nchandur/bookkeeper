package api

import (
	"bookkeeper/internal/db"
	"bookkeeper/internal/recommend"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthorHandler(r *gin.Engine) {
	r.GET("/author", func(ctx *gin.Context) {
		name := ctx.Query("name")

		if name == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"body":  nil,
				"error": "empty name",
			})
			return
		}

		collection := db.Client.Database("booksV2").Collection("works")

		author, err := recommend.GetAuthor(ctx, collection, name)

		fmt.Println(author)

		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"body":  nil,
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"body":  author,
			"error": nil,
		})

	})
}
