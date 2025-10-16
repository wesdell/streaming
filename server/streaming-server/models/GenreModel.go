package models

type Genre struct {
	ID   int    `bson:"genre_id" json:"genre_id" validate:"required"`
	Name string `bson:"genre_name" json:"genre_name" validate:"required,min=3,max=100"`
}
