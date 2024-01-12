package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mirzha99/go-penduduk/config"
	"github.com/mirzha99/go-penduduk/helper"
	"github.com/mirzha99/go-penduduk/models/Muser"
	"github.com/mirzha99/timegoza/timegoza"
	"golang.org/x/crypto/bcrypt"
)

// @Summary  User List
// @Tags User
// @Description Show All Data User
// @Security ApiKeyAuth
// @Accept json
// @Success 200 {object} helper.SuccessResponse
// @Failure 404 {object} helper.ErrorResponse
// @Router /users [get]
func Index(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		e := helper.ErrorResponse{StatusCode: 492, Error: "Too many requests"}
		ctx.JSON(http.StatusTooManyRequests, e.Error)
		ctx.Abort()
		return
	}
	var user []Muser.User

	config.DB.Select("id, nama, email, username, email, role, created_at, change_at").Find(&user)
	if len(user) == 0 {
		e := helper.ErrorResponse{StatusCode: 404, Detail: "Data User Is Empty"}
		ctx.JSON(404, e.Detail)
		return
	}
	for i := range user {
		created_at, _ := strconv.Atoi(user[i].Created_at)
		change_at, _ := strconv.Atoi(user[i].Change_at)
		created := timegoza.ZaTimes{Epoch: int64(created_at), TimeZone: "Asia/Jakarta"}
		change := timegoza.ZaTimes{Epoch: int64(change_at), TimeZone: "Asia/Jakarta"}
		user[i].Created_at = created.HumanTime()
		user[i].Change_at = change.HumanTime()
	}
	tokenCookie, _ := ctx.Cookie("jwt-token")
	s := helper.SuccessResponse{StatusCode: 200, Result: Muser.GetUserAllPublic(user), Token: tokenCookie}
	ctx.JSON(200, s.SuccesResult())
}

// @Summary Get user by id
// @Tags User
// @Description Data User By Id
// @Accept json
// @Param id path int true "User ID"
// @Security ApiKeyAuth
// @Success 201 {object} helper.SuccessResponse
// @Failure 400 {object} helper.ErrorResponse
// @Router /user/{id} [get]
func Byid(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		e := helper.ErrorResponse{StatusCode: 492, Error: "Too many requests"}
		ctx.JSON(http.StatusTooManyRequests, e.Error)
		ctx.Abort()
		return
	}
	var user Muser.User
	id := ctx.Param("id")
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		e := helper.ErrorResponse{StatusCode: 404, Detail: "User Not Found", Error: err.Error()}
		ctx.JSON(404, e.ErrorResultDetail())
		return
	}
	//get cookie
	tokenCookie, err := ctx.Cookie("jwt-token")
	if err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized)
	}
	s := helper.SuccessResponse{Result: user.PublicUser(), StatusCode: 200, Token: tokenCookie}
	ctx.JSON(200, s.SuccesResult())
}
func email_already_exits(email string) bool {
	var user Muser.User
	if err := config.DB.Where("email = ?", email).First(&user).Error; err != nil {
		return false
	}
	return true
}
func username_already_exits(username string) bool {
	var user Muser.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return false
	}
	return true
}

// @Summary Register a new user
// @Tags Auth
// @Description Create a new user with register
// @Security ApiKeyAuth
// @Accept json
// @Param user body Muser.UserInput true "User information"
// @Success 201 {object} helper.SuccessResponse
// @Failure 400 {object} helper.ErrorResponse
// @Router /register [post]
func Add(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		e := helper.ErrorResponse{StatusCode: 492, Error: "Too many requests"}
		ctx.JSON(http.StatusTooManyRequests, e.Error)
		ctx.Abort()
		return
	}
	var user Muser.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(400, gin.H{"error": err.Error()})
		return
	}
	if email_already_exits(user.Email) {
		ctx.JSON(http.StatusConflict, gin.H{"Message": "Email Duplicate!"})
		return
	}
	if username_already_exits(user.Username) {
		ctx.JSON(http.StatusConflict, gin.H{"Message": "Username Duplicate!"})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		e := helper.ErrorResponse{Detail: "Error Hash Password", Error: err.Error(), StatusCode: http.StatusBadRequest}
		ctx.JSON(http.StatusBadRequest, e.ErrorResultDetail())
		return
	}
	user.Password = string(hash)
	user.Created_at = strconv.Itoa(int(timegoza.EpochTime()))
	user.Change_at = strconv.Itoa(int(timegoza.EpochTime()))
	user.Role = "Staff"

	result := config.DB.Create(&user)
	if result.Error != nil {
		ctx.JSON(400, gin.H{"message": "User created Failed", "user": user})
		return
	} else {
		ctx.JSON(201, gin.H{"message": "User created successfully", "user": user.PublicUser()})
	}

}

// @Summary Edit user
// @Tags User
// @Description Edit Data User
// @Security ApiKeyAuth
// @Accept json
// @Param id path int true "User ID"
// @Param user body Muser.UserInput true "User information"
// @Success 201 {object} helper.SuccessResponse
// @Failure 400 {object} helper.ErrorResponse
// @Router /user/{id} [put]
func Edit(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		e := helper.ErrorResponse{StatusCode: 492, Error: "Too many requests"}
		ctx.JSON(http.StatusTooManyRequests, e.Error)
		ctx.Abort()
		return
	}
	var user Muser.User
	id := ctx.Param("id")
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "User not found"})
		return
	}
	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), 10)
	if err != nil {
		e := helper.ErrorResponse{Detail: "Error Hash Password", Error: err.Error(), StatusCode: http.StatusBadRequest}
		ctx.JSON(http.StatusBadRequest, e.ErrorResultDetail())
		return
	}
	user.Password = string(hash)
	user.Change_at = strconv.Itoa(int(timegoza.EpochTime()))
	row := config.DB.Save(&user)
	if row.RowsAffected == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "User update failed"})
		return
	}
	ctx.JSON(200, gin.H{"message": "User Successly Update", "user": user})
}

// @Summary Delete user
// @Tags User
// @Description Delete Data User by Id
// @Security ApiKeyAuth
// @Accept json
// @Param id path int true "User ID"
// @Success 201 {object} helper.SuccessResponse
// @Failure 400 {object} helper.ErrorResponse
// @Router /user/{id} [delete]
func Delete(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		e := helper.ErrorResponse{StatusCode: 492, Error: "Too many requests"}
		ctx.JSON(http.StatusTooManyRequests, e.Error)
		ctx.Abort()
		return
	}
	var user Muser.User
	id := ctx.Param("id")
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "User not found"})
		return
	}
	// Delete the user
	if err := config.DB.Delete(&user).Error; err != nil {
		// Error while deleting user
		ctx.JSON(500, gin.H{"error": "Failed to delete user"})
		return
	}
	ctx.JSON(200, gin.H{"message": "User deleted successfully"})
}
