package api

import (
	"bookkeeper/internal/db"
	"bookkeeper/internal/recommend"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func RecommendByTitleHandler(r *gin.Engine) {
	r.GET("/recommend", func(ctx *gin.Context) {
		title := ctx.Query("title")

		if title == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"body":  nil,
				"error": "empty title",
			})
			return
		}

		n := ctx.Query("n")

		if n == "" {
			ctx.JSON(http.StatusOK, gin.H{
				"body":  nil,
				"error": "missing recommendation number",
			})
			return
		}

		topK, err := strconv.Atoi(n)

		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"body":  nil,
				"error": "invalid top K",
			})
			return
		}

		collection := db.Client.Database("booksV2").Collection("works")
		book, topDocs, err := recommend.GetTopKDocuments(collection, title, topK)

		res := []struct {
			Title  string   `json:"title"`
			Author string   `json:"author"`
			Genres []string `json:"genres"`
			Score  float64  `json:"score"`
		}{}

		for _, t := range topDocs {
			res = append(res, struct {
				Title  string   "json:\"title\""
				Author string   "json:\"author\""
				Genres []string "json:\"genres\""
				Score  float64  "json:\"score\""
			}{
				Title:  t.Doc.Work.Title,
				Author: t.Doc.Work.Author,
				Genres: t.Doc.Work.Genres,
				Score:  t.Score,
			})
		}

		if err != nil {
			ctx.JSON(http.StatusOK, gin.H{
				"body":  nil,
				"error": err.Error(),
			})
			return
			}

		ctx.JSON(http.StatusOK, gin.H{
			"body": gin.H{
				"title":              book.Title,
				"recommended_titles": res,
			},
			"error": nil,
		})

	})
}
