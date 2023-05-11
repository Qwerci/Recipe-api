package models

import (
	"time"

)

type Recipe struct {
	ID			string `json:"id"`
	Name 		string `json:"name"`
	Tags 		[]string `json:"tags"`
	Ingredient 	[]string `json:"ingredient"`
	Instruction []string `json:"instruction"`
	PublishedAt time.Time `json:"published"`
}