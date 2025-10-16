package routes

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/wesdell/streaming/server/streaming-server/controllers"
)

func SetUpUnprotectedRoutes(router *gin.Engine, client *mongo.Client) {
	// User routes
	router.POST("/register", controllers.Register(client))
	router.POST("/login", controllers.Login(client))
	router.POST("/logout", controllers.Logout(client))
	router.POST("/refresh", controllers.RefreshToken(client))

	// Movie routes
	router.GET("/movies", controllers.GetAllMovies(client))

	// Genre routes
	router.GET("/genres", controllers.GetAllGenres(client))
}
