package helpers

import (
	"log"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v4"
)

type SignedDetails struct{
	Email 			string
	First_name 		string
	Last_name 		string
	Uid 			string
	User_type 		string
	jwt.StandardClaims
}

// var userCollection *mongo.Collection = database.OpenCollection(database.Client,"user")

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

	signedToken, err = jwt.NewWithClaims(jwt.SigningMethodES256, claims).SignedString([]byte(SECRET_KEY))
	if err!=nil{
		log.Panic(err.Error())
		return
	}
	signedRefreshToken, err = jwt.NewWithClaims(jwt.SigningMethodES256, refreshClaims).SignedString([]byte(SECRET_KEY))
	if err!=nil{
		log.Panic(err.Error())
		return
	}
	return
}	