package utils

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/wesdell/streaming/server/streaming-server/database"
)

func GetUserIdFromContext(c *gin.Context) (string, error) {
	userId, exists := c.Get("userId")
	if !exists {
		return "", errors.New("user id does not exist")
	}

	id, ok := userId.(string)
	if !ok {
		return "", errors.New("unable to retrieve user id")
	}

	return id, nil
}

func GetUserFavoriteGenres(userId string, client *mongo.Client, c *gin.Context) ([]string, error) {
	ctx, cancel := context.WithTimeout(c, 100*time.Second)
	defer cancel()

	filter := bson.D{{Key: "user_id", Value: userId}}
	projection := bson.M{
		"favorite_genres.genre_name": 1,
		"_id":                        0,
	}
	opt := options.FindOne().SetProjection(projection)

	var result bson.M

	var userCollection = database.OpenCollection("users", client)
	err := userCollection.FindOne(ctx, filter, opt).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []string{}, nil
		}
	}

	favoriteGenres, ok := result["favorite_genres"].(bson.A)
	if !ok {
		return []string{}, errors.New("unable to retrieve favorite genres")
	}

	var genresName []string
	for _, item := range favoriteGenres {
		if genreMap, ok := item.(bson.D); ok {
			for _, elem := range genreMap {
				if elem.Key == "genre_name" {
					if name, ok := elem.Value.(string); ok {
						genresName = append(genresName, name)
					}
				}
			}
		}
	}

	return genresName, nil
}

func GetRoleFromContext(c *gin.Context) (string, error) {
	role, exists := c.Get("role")
	if !exists {
		return "", errors.New("role does not exist")
	}

	memberRole, ok := role.(string)
	if !ok {
		return "", errors.New("unable to retrieve user memberRole")
	}

	return memberRole, nil
}
