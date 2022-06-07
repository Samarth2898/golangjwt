package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Samarth2898/golangjwt/database"
	"github.com/Samarth2898/golangjwt/helpers"
	"github.com/Samarth2898/golangjwt/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)


var userCollection *mongo.Collection = database.OpenCollection(database.Client,"user")
var validate = validator.New()

func HashPassword()


func VerifyPassword()


func Signup() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var dctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		if err := ctx.BindJSON(&user); err!=nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(user)
		if validationErr!=nil{
			ctx.JSON(http.StatusBadRequest, gin.H{"error":validationErr.Error()})
			return
		}

		count, err := userCollection.CountDocuments(dctx, bson.M{"email":user.Email})
		defer cancel()
		if err!=nil{
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while checking for email"})
			return
		}

		phoneCount, err := userCollection.CountDocuments(dctx, bson.M{"phone":user.Phone})
		defer cancel()
		if err!=nil{
			log.Panic(err)
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while checking for phone number"})
			return
		}

		if count>0 || phoneCount>0{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":"this email or phone number already exits"})
			return
		}
	}
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
		var dctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		err := userCollection.FindOne(dctx, bson.M{"user_id":user_id}).Decode(&user)
		defer cancel()
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
		}
		ctx.JSON(http.StatusOK,user)
	}
}