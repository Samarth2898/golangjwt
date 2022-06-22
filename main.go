package main

import (
	"log"
	"net/http"
	"os"

	routes "github.com/Samarth2898/golangjwt/routes"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)


func main(){
	err := godotenv.Load(".env")
	if err!=nil{
		log.Fatal("Error loading the .env file")
	}

	port := os.Getenv("PORT")

	if port == "" {
		port="8000"
	}

	router := gin.New()
	router.Use(gin.Logger())
	routes.AuthRoutes(router)
	routes.UserRoutes(router)
	router.GET("/api-1", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"success":"Access granted for api-1"})
	})
	router.GET("/api-2", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"success":"Access granted for api-2"})
	})

	router.Run(":" + port)
}