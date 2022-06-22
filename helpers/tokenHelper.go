package helpers

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/Samarth2898/golangjwt/database"
	jwt "github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct{
	Email 			string
	First_name 		string
	Last_name 		string
	Uid 			string
	User_type 		string
	jwt.StandardClaims
}

var userCollection *mongo.Collection = database.OpenCollection(database.Client,"user")

var SECRET_KEY = os.Getenv("SECRET_KEY")

func GenerateAllTokens(email string, firstName string, lastName string, userType string, uid string) (signedToken string, signedRefreshToken string, err error){
	
	claims := &SignedDetails{
		Email: email,
		First_name: firstName,
		Last_name: lastName,
		User_type: userType,
		Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}

	refreshClaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(186)).Unix(),
		},
	}
	
	signedToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err!=nil{
		log.Panic(err.Error())
		return
	}
	
	signedRefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err!=nil{
		log.Panic(err.Error())
		return
	}
	return
}	


func ValidateToken(signedToken string) (claims *SignedDetails, msg string){
	token, err := jwt.ParseWithClaims(
		signedToken,
		&SignedDetails{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		},
	)

	if err!=nil{
		msg = err.Error()
		return
	}

	claims, ok := token.Claims.(*SignedDetails)
	if !ok{
		msg = "the token is invalid"
		msg = err.Error()
		return
	}

	if claims.ExpiresAt < int64(time.Now().Local().Unix()){
		msg = "token has expired"
		msg = err.Error()	
	}
	return claims,msg
}



func UpdateAllTokens(signedToken string, signedRefreshToken string, userID string){
	ctx, cancel := context.WithTimeout(context.Background(),100*time.Second)
	var updateObj primitive.D

	updateObj = append(updateObj, bson.E{Key:"token",Value:signedToken})
	updateObj = append(updateObj, bson.E{Key:"refreshToken", Value:signedRefreshToken})

	Updated_at, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateObj = append(updateObj, bson.E{Key:"updated_at", Value:Updated_at})

	upsert := true
	filter := bson.M{"user_id": userID}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}

	_, err := userCollection.UpdateOne(
		ctx,
		filter,
		bson.D{
			{Key:"$set", Value:updateObj},
		},
		&opt,
	) 

	defer cancel()
	if err!=nil{
		log.Panic(err.Error())
		return
	}
	
}