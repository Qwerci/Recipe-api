package controllers

import (
	"os"
	"net/http"
	"time"
	"encoding/json"
	"github.com/Qwerci/Recipe-api/models"
	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
	"strings"
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


func DeleteRecipe(c *gin.Context){
	id := c.Param("id")
	index := -1

	for i := 0; i < len(recipes); i++{
		if recipes[i].ID == id {
			index =1
		}
	}

	if index == -1 {
		c.JSON(http.StatusNotFound,gin.H{
			"error": "Recipes not found",
		})
		return
	}

	recipes = append(recipes[:index], recipes[index+1:]...)
	c.JSON(http.StatusOK, gin.H{
		"message": "Recipes has been deleted",
	})
	
}



func SearchRecipe(c *gin.Context) {

	tag := c.Query("tag")
	listofRecipes := make([]models.Recipe, 0)

	for i := 0; i < len(recipes); i++ {
		found := false

		for _, t := range recipes[i].Tags {
			if strings.EqualFold(t, tag){
				found = true
			}
		}

		if found{
			listofRecipes = append(listofRecipes, recipes[i])

		}
	}

	c.JSON(http.StatusOK, listofRecipes)
}