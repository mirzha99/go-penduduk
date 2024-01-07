package user

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mirzha99/go-penduduk/config"
	"github.com/mirzha99/go-penduduk/models/Muser"
	"github.com/mirzha99/timegoza/timegoza"
	"golang.org/x/crypto/bcrypt"
)

// SuccessResponse adalah struktur respons sukses
type SuccessResponse struct {
	Message string `json:"message"`
	Token   string `json:"token"`
}

// ErrorResponse adalah struktur respons kesalahan
type ErrorResponse struct {
	Error  string `json:"error"`
	Detail string `json:"detail"`
}

func Index(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		ctx.Abort()
		return
	}
	var user []Muser.User

	config.DB.Select("id, nama, email, username, email, role, created_at, change_at").Find(&user)
	if len(user) == 0 {
		ctx.JSON(404, gin.H{"Message": "Data User Is Empty"})
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
	ctx.JSON(200, gin.H{"User": Muser.GetUserAllPublic(user)})
}
func Byid(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		ctx.Abort()
		return
	}
	var user Muser.User
	id := ctx.Param("id")
	if err := config.DB.Where("id = ?", id).First(&user).Error; err != nil {
		ctx.JSON(404, gin.H{"error": "User not found"})
		return
	}

	ctx.JSON(200, gin.H{"User": user.PublicUser()})
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
// @Description Create a new user with the provided data
// @Accept json
// @Param user body Muser.User true "User information"
// @Success 201 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /register [post]
func Add(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
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
		ctx.JSON(http.StatusBadRequest, gin.H{"Message": "Error Hash Password"})
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
func Edit(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
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
	row := config.DB.Save(&user)
	if row.RowsAffected == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "User update failed"})
		return
	}
	ctx.JSON(200, gin.H{"message": "User Successly Update", "user": user})
}
func Delete(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
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
