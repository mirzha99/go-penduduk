package config

import (
	"fmt"
	"os"

	Mkepalakularga "github.com/mirzha99/go-penduduk/models/MKepalaKularga"
	"github.com/mirzha99/go-penduduk/models/Mdesa"
	"github.com/mirzha99/go-penduduk/models/Mmukim"
	"github.com/mirzha99/go-penduduk/models/Muser"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectionDB() {
	defer func() {
		recover := recover()
		if recover != nil {
			fmt.Println(recover)
			os.Exit(0)
		}
	}()
	dns := os.Getenv("dns_db")
	db, err := gorm.Open(mysql.Open(dns), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Muser.User{}, &Mmukim.Mukim{}, &Mdesa.Desa{}, &Mkepalakularga.KepalaKeluarga{})
	DB = db
}
