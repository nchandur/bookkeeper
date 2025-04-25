package api

import (
	"bookkeeper/internal/db"
	"bookkeeper/internal/recommend"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SearchHandler(r *gin.Engine) {
	r.GET("/search", func(ctx *gin.Context) {
		title := ctx.Query("title")

		if title == "" {
			ctx.JSON(http.StatusNotAcceptable, gin.H{
				"body":  nil,
				"error": "empty title",
			})
			return
		}

		collection := db.Client.Database("booksV2").Collection("works")
		doc, err := recommend.GetDocument(collection, title)

		if err != nil {
			ctx.JSON(http.StatusNoContent, gin.H{
				"body":  nil,
				"error": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"body": gin.H{
				"book_id": doc.Work.ID,
				"title":   doc.Work.Title,
				"author":  doc.Work.Author,
				"summary": doc.Work.Summary,
				"genres":  doc.Work.Genres,
				"stars":   doc.Work.Stars,
				"ratings": doc.Work.Ratings,
				"reviews": doc.Work.Reviews,
				"format": gin.H{
					"page_no": doc.Work.Format.PageNo,
					"format":  doc.Work.Format.Type,
				},
				"published": doc.Work.Published,
			},
			"error": nil,
		})

	})
}
