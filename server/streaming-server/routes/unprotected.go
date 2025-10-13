package routes

import (
	"github.com/gin-gonic/gin"

	"github.com/wesdell/streaming/server/streaming-server/controllers"
)

func SetUpUnprotectedRoutes(router *gin.Engine) {
	// User routes
	router.POST("/register", controllers.Register())
	router.POST("/login", controllers.Login())

	//Movie routes
	router.GET("/movies", controllers.GetAllMovies())
}
