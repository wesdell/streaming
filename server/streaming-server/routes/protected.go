package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/wesdell/streaming/server/streaming-server/controllers"
	"github.com/wesdell/streaming/server/streaming-server/middlewares"
)

func SetUpProtectedRoutes(router *gin.Engine) {
	router.Use(middlewares.AuthMiddleware())

	// Movie routes
	router.GET("/movies/:imdb_id", controllers.GetMovieById())
	router.GET("movies/recommended", controllers.GetRecommendedMovies())
	router.POST("/movies", controllers.CreateMovie())
	router.PATCH("/movies/:imdb_id/reviews", controllers.CreateReview())
}
