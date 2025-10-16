package database

import (
	"log"

	"github.com/wesdell/streaming/server/streaming-server/config"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Connect() *mongo.Client {
	mongodb := config.GetEnvVariable("MONGODB_URI")
	if mongodb == "" {
		log.Fatal("Database connection string not set!")
	}

	clientOptions := options.Client().ApplyURI(mongodb)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		return nil
	}
	return client
}

func OpenCollection(collectionName string, client *mongo.Client) *mongo.Collection {
	dbName := config.GetEnvVariable("DB_NAME")
	collection := client.Database(dbName).Collection(collectionName)
	if collection == nil {
		return nil
	}
	return collection
}
