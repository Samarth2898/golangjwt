package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/Samarth2898/golangjwt/database"
	"github.com/Samarth2898/golangjwt/helpers"
	"github.com/Samarth2898/golangjwt/models"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)


var userCollection *mongo.Collection = database.OpenCollection(database.Client,"user")
var validate = validator.New()

func HashPassword()


func VerifyPassword(userPassword string, providedPassword string)(bool, string){
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err!=nil{
		msg = fmt.Sprint("email or password is incorrect")
		check = false
	}
	return check, msg
}


func Signup() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var dctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		if err := ctx.BindJSON(&user); err!=nil{
			defer cancel()
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		validationErr := validate.Struct(user)
		if validationErr!=nil{
			defer cancel()
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

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339)) 
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		token, refreshToken, _ := helpers.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, *&user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNumber, insertError :=userCollection.InsertOne(dctx, user)
		if insertError!=nil{
			msg := fmt.Sprintf("user item was not created")
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":msg})
			return
		}
		defer cancel()
		ctx.JSON(http.StatusOK, gin.H{
			"status":"success",
			"insertion number": resultInsertionNumber,
		})
	}
}

func Login() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		var dctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var foundUser models.User
		if err := ctx.BindJSON(&user); err!=nil{
			defer cancel()
			ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(dctx, bson.M{"email":user.Email}).Decode(&foundUser)
		defer cancel()
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":"email or password is incorrect"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()


	}
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