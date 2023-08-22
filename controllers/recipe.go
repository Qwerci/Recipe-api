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
var err error

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

	recipes := make([]models.Recipe, 0)
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

	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid recipe ID",
		})
		return
	}

	updateFields := bson.D{
		{Key:"name", Value: recipe.Name},
		{Key:"instruction", Value: recipe.Instructions},
		{Key:"ingredient", Value: recipe.Ingredients},
		{Key:"tags", Value: recipe.Tags},
	}

	_, err = recipeCollection.UpdateOne(ctx, 
		bson.M{"_id": objectId},
		bson.D{{Key:"$set", Value: updateFields},
	})



	if err != nil {
		c.JSON(http.StatusInternalServerError, 
		gin.H{"error": err.Error()})
		return
	}

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

	objectID, err := primitive.ObjectIDFromHex(id)
	if  err != nil {
		c.JSON(http.StatusBadRequest, 
		gin.H{
			"error": "Invalid recipe ID",
		})
		return
	}

	for i := 0; i < len(recipes); i++{
		if recipes[i].ID == objectID {
			index = i
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


// SearchRecipe 	godoc
// @Summary			Search Recipe by Tag
// @description		Search for recipes based on a given tag
// @Param			tag query string true "tag to search recipes by"
// @Produce 		application/json
// @tags			recipes
// @Success			200 {object} []models.Recipe "Recipes found successfully"
// @Router          /search-recipes	[get]
func SearchRecipe(c *gin.Context) {

	tag := c.Query("tag")// Get the "tag" query parameter

	// Create a list to store matched recipes
	listofRecipes := make([]models.Recipe, 0)

	// Iterate through each recipe to find matches
	for _, recipe := range recipes {
		found := false

	// Check if the provided tag matches any of the recipe's tags
		for _, t := range recipe.Tags {
			if strings.EqualFold(t, tag){
				found = true
				break		// Once a match is found, exit the inner loop
			}
		}

		if found{
			listofRecipes = append(listofRecipes, recipe)

		}
	}

	c.JSON(http.StatusOK, listofRecipes)
}