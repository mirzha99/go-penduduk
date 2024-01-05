package middleware

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mirzha99/go-penduduk/config"
	"github.com/mirzha99/go-penduduk/models/Muser"
)

func ReqAuth(ctx *gin.Context) {
	//get cookie
	tokenCookie, err := ctx.Cookie("jwt-token")
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	//parse token
	token, err := jwt.Parse(tokenCookie, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(os.Getenv("key_secret")), nil
	})
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		//check exp
		exp := claims["exp"].(float64)
		if float64(time.Now().Unix()) > exp {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}
		//find the user with username cookie
		//get data user_login from cookie jwt-token
		userLogin, ok := claims["user_login"].(map[string]interface{})
		if !ok {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		var user Muser.User
		if err := config.DB.Where("username = ?", userLogin["username"]).First(&user).Error; err != nil {
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		//attach to req
		ctx.Set("user", user.PublicUser())
		//countinue
		ctx.Next()
	} else {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}

}
