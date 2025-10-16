package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/wesdell/streaming/server/streaming-server/controllers"
	"github.com/wesdell/streaming/server/streaming-server/middlewares"
)

func SetUpProtectedRoutes(router *gin.Engine, client *mongo.Client) {
	// Middlewares
	router.Use(middlewares.AuthMiddleware())

	// Movie routes
	router.GET("/movies/:imdb_id", controllers.GetMovieById(client))
	router.GET("movies/recommended", controllers.GetRecommendedMovies(client))
	router.POST("/movies", controllers.CreateMovie(client))
	router.PATCH("/movies/:imdb_id/reviews", controllers.CreateReview(client))
}
