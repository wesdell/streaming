package models

type Ranking struct {
	Value int    `bson:"value" json:"value" validate:"required"`
	Name  string `bson:"name" json:"name" validate:"required"`
}
