package models

import (
	"time"

)

type Recipe struct {
	Name 		string `json:"name"`
	Tags 		[]string `json:"tags"`
	Ingredient 	[]string `json:"ingredient"`
	Instruction []string `json:"instruction"`
	PublishedAt time.Time `json:"published"`
}