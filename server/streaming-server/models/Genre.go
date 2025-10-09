package models

type Genre struct {
	ID   int    `bson:"_id" json:"_id" validate:"required"`
	Name string `bson:"name" json:"name" validate:"required,min=3,max=100"`
}
