package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Recipe struct {
	ID			primitive.ObjectID `bson:"_id"`
	Name 		string `json:"name"`
	Tags 		[]string `json:"tags"`
	Ingredient 	[]string `json:"ingredient"`
	Instruction []string `json:"instruction"`
	PublishedAt time.Time `json:"publishedAt"`
}