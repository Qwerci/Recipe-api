package main

import (
	"github.com/gin-gonic/gin"
	"github.com/Qwerci/Recipe-api/controllers"
)

func main() {
	router :=gin.Default()

	router.POST("/recipes",controllers.NewRecipe)

	router.Run(":5000")


}