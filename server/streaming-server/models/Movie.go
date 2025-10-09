package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
)

type Movie struct {
	ID          bson.ObjectID `bson:"_id,omitempty" json:"_id,omitempty"`
	ImdbID      string        `bson:"imdb_id" json:"imdb_id" validate:"required"`
	Title       string        `bson:"title" json:"title" validate:"required,min=3,max=150"`
	PosterPath  string        `bson:"poster_path" json:"poster_path" validate:"required,url"`
	YouTubeID   string        `bson:"youtube_id" json:"youtube_id" validate:"required"`
	AdminReview string        `bson:"admin_review" json:"admin_review"`
	Ranking     Ranking       `bson:"ranking" json:"ranking" validate:"required"`
	Genre       []Genre       `bson:"genre" json:"genre" validate:"required,dive"`
}

type Genre struct {
	ID   int    `bson:"_id" json:"_id" validate:"required"`
	Name string `bson:"name" json:"name" validate:"required,min=3,max=100"`
}

type Ranking struct {
	Value int    `bson:"value" json:"value" validate:"required"`
	Name  string `bson:"name" json:"name" validate:"required"`
}
