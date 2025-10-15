package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"

	"github.com/wesdell/streaming/server/streaming-server/database"
	"github.com/wesdell/streaming/server/streaming-server/models"
	"github.com/wesdell/streaming/server/streaming-server/openai"
)

var movieCollection = database.OpenCollection("movies")

var validate = validator.New()

func GetAllMovies() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movies []models.Movie
		cursor, err := movieCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch movies"})
		}
		defer cursor.Close(ctx)

		if err = cursor.All(ctx, &movies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode movies"})
		}

		c.JSON(http.StatusOK, movies)
	}
}

func GetMovieById() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		movieImdbId := c.Param("imdb_id")
		if movieImdbId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie IMDB id is required"})
			return
		}

		var movie models.Movie
		err := movieCollection.FindOne(ctx, bson.M{"imdb_id": movieImdbId}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func CreateMovie() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movie models.Movie
		if err := c.ShouldBindJSON(&movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
			return
		}

		if err := validate.Struct(movie); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func CreateReview() gin.HandlerFunc {
	return func(c *gin.Context) {
		movieId := c.Param("imdb_id")
		if movieId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie IMDB id is required"})
			return
		}

		var req struct {
			AdminReview string `json:"admin_review"`
		}

		var res struct {
			RankingName string `json:"ranking_name"`
			AdminReview string `json:"admin_review"`
		}

		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		sentiment, rankingValue, err := openai.GetReviewRanking(req.AdminReview)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get review ranking"})
			return
		}

		filter := bson.M{"imdb_id": movieId}
		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"ranking_name":  sentiment,
					"ranking_value": rankingValue,
				},
			},
		}

		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		result, err := movieCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed updating movie"})
			return
		}
		if result.MatchedCount == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		res.RankingName = sentiment
		res.AdminReview = req.AdminReview
		c.JSON(http.StatusOK, res)
	}
}
