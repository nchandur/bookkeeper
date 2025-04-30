package api

import "github.com/gin-gonic/gin"

func SetUpRouter() *gin.Engine {
	r := gin.Default()

	LandingHandler(r)
	PingHandler(r)
	SearchHandler(r)
	RecommendByTitleHandler(r)
	AuthorHandler(r)

	return r

}
