package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"golang.org/x/crypto/bcrypt"

	"github.com/wesdell/streaming/server/streaming-server/database"
	"github.com/wesdell/streaming/server/streaming-server/models"
	"github.com/wesdell/streaming/server/streaming-server/utils"
)

var userCollection *mongo.Collection = database.OpenCollection("users")

func Register() gin.HandlerFunc {
	return func(c *gin.Context) {
		var user models.User
		if err := c.ShouldBindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		validate := validator.New()
		if err := validate.Struct(user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedPassword, err := hashPassword(user.Password)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to hash password"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing user"})
			return
		}
		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User already exists"})
			return
		}

		user.UserId = bson.NewObjectID().Hex()
		user.Password = hashedPassword
		user.CreatedAt = time.Now()
		user.UpdatedAt = time.Now()
		result, err := userCollection.InsertOne(ctx, user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"user": result})
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var userLogin models.UserLogin
		if err := c.ShouldBindJSON(&userLogin); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
			return
		}

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var existingUser models.User
		if err := userCollection.FindOne(ctx, bson.M{"email": userLogin.Email}).Decode(&existingUser); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}
		if err := bcrypt.CompareHashAndPassword([]byte(existingUser.Password), []byte(userLogin.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
			return
		}

		token, refreshToken, err := utils.GenerateTokens(existingUser.Email, existingUser.FirstName, existingUser.LastName, existingUser.Role, existingUser.UserId)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
			return
		}

		err = utils.UpdateTokens(existingUser.UserId, token, refreshToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update tokens"})
			return
		}

		c.JSON(http.StatusOK, models.UserResponse{
			UserId:         existingUser.UserId,
			FirstName:      existingUser.FirstName,
			LastName:       existingUser.LastName,
			Role:           existingUser.Role,
			Email:          existingUser.Email,
			FavoriteGenres: existingUser.FavoriteGenres,
			Token:          token,
			RefreshToken:   refreshToken,
		})
	}
}

func hashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}
