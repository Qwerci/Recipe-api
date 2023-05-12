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

func UpdateRecipe(c *gin.Context){
	id := c.Param("id")

	var recipe models.Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil {
		c.JSON(http.StatusOK, gin.H{
			"error": err.Error()})
		return
	}

	index := -1

	for i := 0; i < len(recipes); i++ {
		if recipes[i].ID == id {
			index = i
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound,gin.H{
			"error": "Recipe not found",
		})
		return
	}

	recipes[index] = recipe

	c.JSON(http.StatusOK, recipe)

}




