package middleware

import (
	"net/http"

	helper "github.com/Samarth2898/golangjwt/helpers"
	"github.com/gin-gonic/gin"
)


func Authenticate() gin.HandlerFunc{
	return func(ctx *gin.Context) {
		clientToken := ctx.Request.Header.Get("token")
		if clientToken==""{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":"No authorization header provided."})
			ctx.Abort()
			return
		}

		claims, err := helper.ValidateToken(clientToken)
		if err!=""{
			ctx.JSON(http.StatusInternalServerError, gin.H{"error":err})
			ctx.Abort()
			return
		}

		ctx.Set("email", claims.Email)
		ctx.Set("first_name", claims.First_name)
		ctx.Set("last_name", claims.Last_name)
		ctx.Set("uid",claims.Uid)
		ctx.Set("user_type", claims.User_type)
		ctx.Next()	
	}
}