package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"github.com/wesdell/streaming/server/streaming-server/config"
	"github.com/wesdell/streaming/server/streaming-server/database"
	"github.com/wesdell/streaming/server/streaming-server/routes"
)

func main() {
	router := gin.Default()

	allowedOrigins := config.GetEnvVariable("ALLOWED_ORIGINS")

	var origins []string
	if allowedOrigins != "" {
		origins = strings.Split(allowedOrigins, ",")
		for i := range origins {
			origins[i] = strings.TrimSpace(origins[i])
			log.Println("Allowed Origin:", origins[i])
		}
	} else {
		origins = []string{"http://localhost:5173"}
		log.Println("Allowed Origin: http://localhost:5173")
	}

	corsConfig := cors.Config{}
	corsConfig.AllowOrigins = origins
	corsConfig.AllowMethods = []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"}
	//corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Accept", "Authorization"}
	corsConfig.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	corsConfig.ExposeHeaders = []string{"Content-Length"}
	corsConfig.AllowCredentials = true
	corsConfig.MaxAge = 12 * time.Hour

	router.Use(cors.New(corsConfig))
	router.Use(gin.Logger())

	var client = database.Connect()

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatalf("Failed to reach server: %v", err)
	}
	defer func() {
		err := client.Disconnect(context.Background())
		if err != nil {
			log.Fatalf("Failed to disconnect from MongoDB: %v", err)
		}

	}()

	routes.SetUpUnprotectedRoutes(router, client)
	routes.SetUpProtectedRoutes(router, client)

	if err := router.Run(":8080"); err != nil {
		fmt.Println("Failed to start server!", err)
	}
}
