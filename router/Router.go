package router

import (
	"github.com/gin-gonic/gin"
	"github.com/mirzha99/go-penduduk/controller/auth"
	"github.com/mirzha99/go-penduduk/controller/desa"
	"github.com/mirzha99/go-penduduk/controller/home"
	"github.com/mirzha99/go-penduduk/controller/mukim"
	"github.com/mirzha99/go-penduduk/controller/user"
	_ "github.com/mirzha99/go-penduduk/docs"
	"github.com/mirzha99/go-penduduk/middleware"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func Router() {
	r := gin.Default()

	r.GET("/", home.Index)
	//auth
	r.POST("/login", auth.Login)
	r.POST("/register", user.Add)
	r.GET("/profil", middleware.ReqAuth, auth.Profil)
	//user
	r.GET("/users", middleware.ReqAuth, user.Index)
	r.GET("/user/:id", middleware.ReqAuth, user.Byid)
	r.PUT("/user/:id", middleware.ReqAuth, user.Edit)
	r.DELETE("/user/:id", middleware.ReqAuth, user.Delete)
	//mukim
	r.GET("/mukims", middleware.ReqAuth, mukim.Index)
	r.GET("/mukim/:id", middleware.ReqAuth, mukim.GetById)
	r.POST("/mukim/", middleware.ReqAuth, mukim.Add)
	r.PUT("/mukim/:id", middleware.ReqAuth, mukim.Edit)
	r.DELETE("/mukim/:id", middleware.ReqAuth, mukim.Delete)
	//desa
	r.GET("/desas", middleware.ReqAuth, desa.Index)
	r.GET("/desa/:id", middleware.ReqAuth, desa.GetById)
	r.POST("/desa/", middleware.ReqAuth, desa.Add)
	r.PUT("/desa/:id", middleware.ReqAuth, desa.Edit)
	r.DELETE("/desa/:id", middleware.ReqAuth, desa.Delete)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	r.Run(":8080")
}
