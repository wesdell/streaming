package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
	controller "github.com/wesdell/streaming/server/streaming-server/controllers"
)

func main() {
	router := gin.Default()

	router.GET("/movies", controller.GetAllMovies())

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server!", err)
	}
}
