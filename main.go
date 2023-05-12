package main

import (
	"github.com/gin-gonic/gin"
	"github.com/Qwerci/Recipe-api/controllers"
)

func main() {
	router :=gin.Default()

	router.POST("/recipes",controllers.NewRecipe)
	router.GET("/recipes",controllers.ListRecipes)
	router.PUT("/recipes/:id",controllers.UpdateRecipe)
	router.DELETE("/recipes/:id",controllers.DeleteRecipe)

	router.Run(":5000")


}