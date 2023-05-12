package controllers

import (
	"os"
	"net/http"
	"time"
	"encoding/json"
	"github.com/Qwerci/Recipe-api/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)


var recipes []models.Recipe


func init() {
	recipes = make([]models.Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)
}

func NewRecipe(c *gin.Context){
	var recipe models.Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return 
	}

	recipe.ID = xid.New().String()
	recipe.PublishedAt = time.Now()
	recipes = append(recipes, recipe)
	c.JSON(http.StatusOK, recipe)

}


func ListRecipes(c *gin.Context){
	c.JSON(http.StatusOK, recipes)
}

