package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/Qwerci/Recipe-api/db"
	"github.com/Qwerci/Recipe-api/models"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


var recipes []models.Recipe
var recipeCollection *mongo.Collection = db.OpenCollection(db.Client,"recipes")
var ctx = context.Background()

func init() {
	

	recipes = make([]models.Recipe, 0)
	file, _ := os.ReadFile("recipes.json")
	_ = json.Unmarshal([]byte(file), &recipes)

	var listofRecipes []interface{}
	for _, recipe := range recipes {
		listofRecipes = append(listofRecipes, recipe)
	}
	
	insertManyResult, err := recipeCollection.InsertMany(ctx, listofRecipes)
	if err != nil {
		log.Fatal(err)
	}
	
	log.Println("Inserted recipes: ", len(insertManyResult.InsertedIDs))
}

// CreateRecipe 	godoc
// @Summary			Create Recipe
// @description		Create new recipes and saves it to a json file
// 
// @Produce 		application/json
// @tags			recipes
// @Success			200 {object} models.Recipe "Recipe created successfully"
// @Failure			400 {string} string "Bad request"
// @Router          /recipes	[post]
func NewRecipe(c *gin.Context){
	var recipe models.Recipe

	if err := c.ShouldBindJSON(&recipe); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error()})
		return 
	}

	recipe.ID = primitive.NewObjectID()
	recipe.PublishedAt = time.Now()

	var err error
	_, err = recipeCollection.InsertOne(ctx, recipe)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError,
		gin.H{"error": "error creating recipe"})
		return 
	}
	c.JSON(http.StatusOK, recipe)

}
// ListRecipe 	godoc
// @Summary			List Recipe
// @description		List all recipes
// 
// @Produce 		application/json
// @tags			recipes
// @Success			200 {object} models.Recipe "Recipe listed successfully"
// @Router          /recipes	[get]
func ListRecipes(c *gin.Context){
	
	cur, err := recipeCollection.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, 
		gin.H{"error": err.Error()})
		return
	}
	defer cur.Close(ctx)

	recipes := make([] models.Recipe, 0)
	for cur.Next(ctx) {
		var recipe models.Recipe
		cur.Decode(&recipe)
		recipes = append(recipes, recipe)
	}
	c.JSON(http.StatusOK, recipes)
}

// UpdateRecipe 	godoc
// @Summary			Update Recipe
// @description		Update an exiting recipe
// @Param			recipeid path string true "update recipe by id"
// @Produce 		application/json
// @tags			recipes
// @Success			200 {object} models.Recipe "Recipe updated successfully"
// @Failure			404 {string} string "Recipe not found"
// @Router          /recipes/:id	[put]
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

	c.JSON(http.StatusOK, gin.H{"message": "Recipe has been updated"})

}

// UpdateRecipe 	godoc
// @Summary			Delete Recipe
// @description		Delete an exiting recipe
// @Param			recipeid path string true "delete recipe by id"
// @Produce 		application/json
// @tags			recipes
// @Success			200 {object} models.Recipe "Recipe deleted successfully"
// @Failure			404 {string} string "Recipe not found"
// @Router          /recipes/:id	[delete]
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