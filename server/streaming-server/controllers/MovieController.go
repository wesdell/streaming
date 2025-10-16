package controllers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/wesdell/streaming/server/streaming-server/config"
	"github.com/wesdell/streaming/server/streaming-server/database"
	"github.com/wesdell/streaming/server/streaming-server/models"
	"github.com/wesdell/streaming/server/streaming-server/openai"
	"github.com/wesdell/streaming/server/streaming-server/utils"
)

var validate = validator.New()

func GetAllMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieCollection = database.OpenCollection("movies", client)

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

func GetMovieById(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		movieImdbId := c.Param("imdb_id")
		if movieImdbId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Movie IMDB id is required"})
			return
		}

		var movieCollection = database.OpenCollection("movies", client)
		var movie models.Movie
		err := movieCollection.FindOne(ctx, bson.D{{Key: "imdb_id", Value: movieImdbId}}).Decode(&movie)
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
			return
		}

		c.JSON(http.StatusOK, movie)
	}
}

func CreateMovie(client *mongo.Client) gin.HandlerFunc {
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

		var movieCollection = database.OpenCollection("movies", client)
		result, err := movieCollection.InsertOne(ctx, movie)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
			return
		}

		c.JSON(http.StatusCreated, result)
	}
}

func CreateReview(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, err := utils.GetRoleFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to fetch role"})
			return
		}
		if role != "ADMIN" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User must be administrator"})
			return
		}

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

		sentiment, rankingValue, err := openai.GetReviewRanking(req.AdminReview, client, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not get review ranking"})
			return
		}

		filter := bson.D{{Key: "imdb_id", Value: movieId}}
		update := bson.M{
			"$set": bson.M{
				"admin_review": req.AdminReview,
				"ranking": bson.M{
					"name":  sentiment,
					"value": rankingValue,
				},
			},
		}

		ctx, cancel := context.WithTimeout(c, 100*time.Second)
		defer cancel()

		var movieCollection = database.OpenCollection("movies", client)
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

func GetRecommendedMovies(client *mongo.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		userId, err := utils.GetUserIdFromContext(c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get user id"})
			return
		}

		favoriteGenres, err := utils.GetUserFavoriteGenres(userId, client, c)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		var recommendedMoviesLimit int64
		recommendedMoviesLimitStr := config.GetEnvVariable("RECOMMENDED_MOVIE_LIMIT")
		if recommendedMoviesLimitStr != "" {
			recommendedMoviesLimit, _ = strconv.ParseInt(recommendedMoviesLimitStr, 10, 64)
		}

		findOptions := options.Find()
		findOptions.SetSort(bson.D{{Key: "ranking.value", Value: 1}})
		findOptions.SetLimit(recommendedMoviesLimit)

		filter := bson.D{{
			Key: "genre.genre_name", Value: bson.D{{
				Key: "$in", Value: favoriteGenres,
			}},
		}}

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var movieCollection = database.OpenCollection("movies", client)
		cursor, err := movieCollection.Find(ctx, filter, findOptions)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching recommended movies"})
			return
		}
		defer cursor.Close(ctx)

		var recommendedMovies []models.Movie

		if err := cursor.All(ctx, &recommendedMovies); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, recommendedMovies)

	}
}
