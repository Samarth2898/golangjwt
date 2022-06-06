package controllers

import (
	"net/http"

	"github.com/Samarth2898/golangjwt/database"
	"github.com/Samarth2898/golangjwt/helpers"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)


var userCollection *mongo.Collection = database.OpenCollection(database.Client,"user")
var validate = validator.New()

func HashPassword()


func VerifyPassword()


func Signup(){

}

func Login(){

}

func GetUser(){

}

func GetUsers() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		user_id := ctx.Param("user_id")
		if err:= helpers.MatchUserTypeToUid(ctx, user_id); err!=nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
	}
}