package auth

import (
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mirzha99/go-penduduk/config"
	"github.com/mirzha99/go-penduduk/helper"
	"github.com/mirzha99/go-penduduk/models/Muser"
	"github.com/mirzha99/timegoza/timegoza"
	"golang.org/x/crypto/bcrypt"
)

// @Summary Login
// @Tags Auth
// @Description Login with the provided data
// @Accept json
// @Param user body Muser.LoginInput true "User information"
// @Success 201 {object} helper.SuccessResponse
// @Failure 400 {object} helper.ErrorResponse
// @Router /login [post]
func Login(ctx *gin.Context) {
	var loginInput Muser.LoginInput
	var user Muser.User
	if err := ctx.ShouldBindJSON(&loginInput); err != nil {
		e := helper.ErrorResponse{Error: "JSON Invalid", Detail: err.Error()}
		ctx.JSON(400, gin.H{"error": e.Error, "detail": e.Detail})
		return
	}

	if err := config.DB.Where("username = ?", loginInput.Username).Find(&user).Error; err != nil {
		e := helper.ErrorResponse{Error: "Username Not Found", Detail: err.Error()}
		ctx.JSON(400, gin.H{"error": e.Error, "detail": e.Detail})
		return
	}
	//compare hash password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginInput.Password))
	if err != nil {
		e := helper.ErrorResponse{Error: "Password Invalid", Detail: err.Error()}
		ctx.JSON(400, gin.H{"error": e.Error, "detail": e.Detail})
		return
	}
	//jwt
	created_at, _ := strconv.Atoi(user.Created_at)
	change_at, _ := strconv.Atoi(user.Change_at)
	created := timegoza.ZaTimes{Epoch: int64(created_at), TimeZone: "Asia/Jakarta"}
	change := timegoza.ZaTimes{Epoch: int64(change_at), TimeZone: "Asia/Jakarta"}
	user.Created_at = created.HumanTime()
	user.Change_at = change.HumanTime()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_login": user.PublicUser(),
		"exp":        time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	tokenResult, err := token.SignedString([]byte(os.Getenv("key_secret")))
	if err != nil {
		e := helper.ErrorResponse{Error: "Token Fail Cretaed", Detail: err.Error()}
		ctx.JSON(400, e.ErrorResultDetail())
		return
	}

	//json
	ctx.SetSameSite(http.SameSiteLaxMode)
	ctx.SetCookie("jwt-token", tokenResult, 3600*24*30, "", "", false, true)
	s := helper.SuccessResponse{StatusCode: 200, Message: "Welcome " + user.Nama, Token: tokenResult}
	ctx.JSON(200, s.SuccesRMessage())

}

// @Summary Logout
// @Tags Auth
// @Description Logout user
// @Accept json
// @Success 200
// @Router /logout [get]
func Logout(ctx *gin.Context) {
	ctx.SetCookie("jwt-token", "", -1, "", "", false, true)
	s := helper.SuccessResponse{StatusCode: 200, Message: "Logout Berhasil"}
	ctx.JSON(200, s.SuccesRMessage())
}

// @Summary Profil User Login
// @Tags Auth
// @Security ApiKeyAuth
// @Description Login with the provided data
// @Accept json
// @Success 201 {object} helper.SuccessResponse
// @Failure 400 {object} helper.ErrorResponse
// @Router /profil [get]
func Profil(ctx *gin.Context) {
	profil, err := ctx.Get("user")
	if !err {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	s := helper.SuccessResponse{StatusCode: 200, Result: profil}
	ctx.JSON(200, s.Result)
}
