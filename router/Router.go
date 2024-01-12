package router

import (
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/mirzha99/go-penduduk/controller/auth"
	datakk "github.com/mirzha99/go-penduduk/controller/data-kk"
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
	// Konfigurasi CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost/" + os.Getenv("PORT")} // Ganti dengan alamat domain Anda
	r.Use(cors.New(config))

	// Menentukan endpoint untuk file statis
	r.Static("/datapenduduk", "./uploads/datapenduduk")

	r.GET("/", home.Index)

	//auth
	r.POST("/login", auth.Login)
	r.POST("/register", user.Add)
	r.GET("/profil", middleware.ReqAuth, auth.Profil)
	r.GET("/logout", auth.Logout)

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

	//data kepala keluarga
	r.GET("/kepalakeluarga", datakk.Index)
	r.GET("/kepalakeluarga/:id", datakk.ById)
	r.POST("/kepalakeluarga", datakk.Add)

	//swagger index
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	//run framework gin by .env varibel port
	r.Run(os.Getenv("PORT"))
}
