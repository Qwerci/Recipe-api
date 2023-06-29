// @title		Recipes API
// @description	This is a sample recipes API implementation. It's a learning jounery to understand
// the tools and packages in Go

// @schemes http
// @host localhost:5000
// @BasePath	/
// @version	Version: 1.0.0
// 

package main

import (
	_ "github.com/Qwerci/Recipe-api/docs"
	"github.com/gin-gonic/gin"
	"github.com/Qwerci/Recipe-api/controllers"
	// "github.com/Qwerci/Recipe-api/db"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
)


func main() {
	
	router :=gin.Default()
// add swagger
	router.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/recipes",controllers.NewRecipe)
	router.GET("/recipes",controllers.ListRecipes)
	router.PUT("/recipes/:id",controllers.UpdateRecipe)
	router.DELETE("/recipes/:id",controllers.DeleteRecipe)
	router.GET("/recipes/search",controllers.SearchRecipe)

	router.Run(":5000")


}