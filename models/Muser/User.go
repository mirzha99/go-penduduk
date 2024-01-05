package Muser

import "github.com/gin-gonic/gin"

type User struct {
	Id         int    `gorm:"PrimaryKey" column:"id" json:"id"`
	Nama       string `gorm:"type:varchar(30);column:nama" json:"nama"`
	Email      string `gorm:"type:varchar(30);column:email" json:"email"`
	Username   string `gorm:"type:varchar(30);column:username" json:"username"`
	Password   string `gorm:"type:varchar(60);column:password" json:"password"`
	Role       string `gorm:"type:varchar(60);column:role" json:"role"`
	Created_at string `gorm:"type:varchar(30);column:created_at" json:"created_at"`
	Change_at  string `gorm:"type:varchar(30);column:change_at" json:"change_at"`
}
type LoginInput struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (u *User) PublicUser() gin.H {
	return gin.H{
		"id":         u.Id,
		"nama":       u.Nama,
		"email":      u.Email,
		"username":   u.Username,
		"role":       u.Role,
		"created_at": u.Created_at,
		"change_at":  u.Change_at,
	}
}
func GetUserAllPublic(user []User) []gin.H {
	var publicUser []gin.H
	for _, user := range user {
		publicUser = append(publicUser, user.PublicUser())
	}
	return publicUser
}
