package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/wesdell/streaming/server/streaming-server/controllers"
)

func main() {
	router := gin.Default()

	router.GET("/movies", controllers.GetAllMovies())
	router.GET("/movies/:imdb_id", controllers.GetMovieById())

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server!", err)
	}
}
