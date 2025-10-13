package main

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/wesdell/streaming/server/streaming-server/routes"
)

func main() {
	router := gin.Default()

	routes.SetUpUnprotectedRoutes(router)
	routes.SetUpProtectedRoutes(router)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server!", err)
	}
}
