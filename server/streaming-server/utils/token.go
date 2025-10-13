package utils

import (
	"context"
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/wesdell/streaming/server/streaming-server/config"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/wesdell/streaming/server/streaming-server/database"
)

type Token struct {
	FirstName string
	LastName  string
	Email     string
	Role      string
	UserId    string
	jwt.RegisteredClaims
}

var secretKey string = config.GetEnvVariable("JWT_SECRET_KEY")
var refreshKey string = config.GetEnvVariable("JWT_SECRET_REFRESH_KEY")

func GenerateTokens(email, firstName, lastName, role, userId string) (string, string, error) {
	claims := &Token{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Streaming",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", "", err
	}

	refreshClaims := &Token{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Role:      role,
		UserId:    userId,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "Streaming",
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
		},
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefreshToken, err := refreshToken.SignedString([]byte(refreshKey))
	if err != nil {
		return "", "", err
	}

	return signedToken, signedRefreshToken, nil
}

func UpdateTokens(userId, token, refreshToken string) (err error) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	updateAt, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	update := bson.M{
		"$set": bson.M{
			"token":         token,
			"refresh_token": refreshToken,
			"updated_at":    updateAt,
		},
	}

	var userCollection *mongo.Collection = database.OpenCollection("users")
	_, err = userCollection.UpdateOne(ctx, bson.M{"user_id": userId}, update)
	if err != nil {
		return err
	}
	return nil
}

func GetToken(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", errors.New("No Authorization header found")
	}

	tokenString := authHeader[len("Bearer "):]
	if tokenString == "" {
		return "", errors.New("No token found")
	}

	return tokenString, nil
}
