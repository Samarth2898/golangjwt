package controllers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/Samarth2898/golangjwt/database"
	helper "github.com/Samarth2898/golangjwt/helpers"
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

func HashPassword(password string) string{
	pass, err := bcrypt.GenerateFromPassword([]byte(password),14)
	if err!=nil{
		log.Panic(err.Error())
	}
	return string(pass)
}


func VerifyPassword(userPassword string, providedPassword string)(bool, string){
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err!=nil{
		msg = "email or password is incorrect"
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
		password := HashPassword(*user.Password)
		user.Password = &password


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
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, *user.User_type, user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNumber, insertError :=userCollection.InsertOne(dctx, user)
		if insertError!=nil{
			msg := "user item was not created"
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
		if !passwordIsValid {
			ctx.JSON(http.StatusInternalServerError,gin.H{"error":msg})
			return
		}

		if foundUser.Email == nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":"user not found"})
		}

		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email,*foundUser.First_name,*foundUser.Last_name,*foundUser.User_type,foundUser.User_id)
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		err = userCollection.FindOne(dctx, bson.M{"user_id": foundUser.User_id}).Decode(&foundUser)
		if err!=nil{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":err.Error()})
			return
		}

		ctx.JSON(http.StatusOK, foundUser)
	}
}

func GetUsers() gin.HandlerFunc{
	return func(c *gin.Context){
		if err := helper.CheckUserType(c, "ADMIN"); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error":err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		
		recordPerPage, err := strconv.Atoi(c.Query("recordPerPage"))
		if err != nil || recordPerPage <1{
			recordPerPage = 10
		}
		page, err1 := strconv.Atoi(c.Query("page"))
		if err1 !=nil || page<1{
			page = 1
		}

		startIndex := (page - 1) * recordPerPage
		startIndex, err = strconv.Atoi(c.Query("startIndex"))

		matchStage := bson.D{primitive.E{Key:"$match", Value:bson.D{{}}}}
		groupStage := bson.D{primitive.E{Key:"$group",Value:bson.D{
			primitive.E{Key:"_id", Value:bson.D{primitive.E{Key:"_id", Value: "null"}}}, 
			primitive.E{Key:"total_count", Value: bson.D{primitive.E{Key:"$sum", Value: 1}}}, 
			primitive.E{Key:"data", Value: bson.D{primitive.E{Key:"$push", Value: "$$ROOT"}}}}}}
		projectStage := bson.D{
			{Key:"$project",Value:  bson.D{
				{Key:"_id", Value:  0},
				{Key:"total_count", Value: 1},
				{Key:"user_items", Value:bson.D{{Key:"$slice", Value:[]interface{}{"$data", startIndex, recordPerPage}}}},}}}
		result,err := userCollection.Aggregate(ctx, mongo.Pipeline{
			matchStage, groupStage, projectStage})
		defer cancel()
		if err!=nil{
			c.JSON(http.StatusInternalServerError, gin.H{"error":"error occured while listing user items"})
		}
		var allusers []bson.M
		if err = result.All(ctx, &allusers); err!=nil{
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, allusers[0])
}}

func GetUser() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		user_id := ctx.Param("user_id")
		if err:= helper.MatchUserTypeToUid(ctx, user_id); err!=nil{
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