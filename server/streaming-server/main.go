package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/wesdell/streaming/server/streaming-server/controllers"
)

func main() {
	router := gin.Default()

	// User routes
	router.POST("/register", controllers.Register())
	router.POST("/login", controllers.Login())

	// Movie routes
	router.GET("/movies", controllers.GetAllMovies())
	router.GET("/movies/:imdb_id", controllers.GetMovieById())
	router.POST("/movies", controllers.CreateMovie())

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server!", err)
	}
}
