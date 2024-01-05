package desa

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mirzha99/go-penduduk/config"
	"github.com/mirzha99/go-penduduk/models/Mdesa"
	"github.com/mirzha99/go-penduduk/models/Mmukim"
)

//var desa Mdesa.Desa

func Index(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		ctx.Abort()
		return
	}
	var desaresults []Mdesa.DesaResult
	config.DB.Table("desas").
		Select("desas.id, desas.nama, desas.nama_kepala_desa, desas.id_mukim, mukims.nama AS nama_mukim, mukims.nama_imum_mukim").
		Joins("JOIN mukims ON mukims.id = desas.id_mukim").
		Find(&desaresults)

	if len(desaresults) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Desa is Empty"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"desa": desaresults})
}
func GetById(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		ctx.Abort()
		return
	}
	var desaresult Mdesa.DesaResult
	id := ctx.Param("id")

	// Menggunakan Find bukan First karena kita ingin mengambil semua kolom yang dipilih
	result := config.DB.Table("desas").
		Select("desas.id, desas.nama, desas.nama_kepala_desa, desas.id_mukim, mukims.nama AS nama_mukim, mukims.nama_imum_mukim").
		Joins("JOIN mukims ON mukims.id = desas.id_mukim").
		Where("desas.id = ?", id).Find(&desaresult)
	// Periksa apakah record ditemukan
	if result.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Desa dengan ID " + id + " tidak ditemukan"})
		return
	} else if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Gagal mengambil data desa"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"desa": desaresult})
}
func desa_already_exits(nama_desa string) bool {
	var desa Mdesa.Desa
	if err := config.DB.Where("nama = ?", nama_desa).First(&desa).Error; err != nil {
		return false
	}
	return true
}
func Add(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		ctx.Abort()
		return
	}
	var desa Mdesa.Desa
	if err := ctx.ShouldBindJSON(&desa); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
	}
	var mukim Mmukim.Mukim
	if row := config.DB.Where("id = ?", desa.Id_mukim).First(&mukim).RowsAffected; row == 0 {
		msg := fmt.Sprintf("Id mukim %v tidak ada !", desa.Id_mukim)
		ctx.JSON(http.StatusNotFound, gin.H{"message": msg})
		return
	}
	if desa_already_exits(desa.Nama) {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Nama Desa Duplikat"})
		return
	}
	row := config.DB.Create(&desa)
	if row.RowsAffected == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "data desa cretad failed!"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "data desa cretad success!", "desa": &desa})
}
func Edit(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		ctx.Abort()
		return
	}
	var desa Mdesa.Desa

	id := ctx.Param("id")
	if row := config.DB.Where("id = ?", id).First(&desa).RowsAffected; row == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Data desa dengan id " + id + " tidak ada !"})
		return
	}

	// if desa_already_exits(desa.Nama) {
	// 	ctx.JSON(http.StatusNotFound, gin.H{"message": "Nama Desa Duplikat"})
	// 	return
	// }
	if err := ctx.ShouldBindJSON(&desa); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	var mukim Mmukim.Mukim
	if row := config.DB.Where("id = ?", desa.Id_mukim).First(&mukim).RowsAffected; row == 0 {
		msg := fmt.Sprintf("Id mukim %v tidak ada !", desa.Id_mukim)
		ctx.JSON(http.StatusNotFound, gin.H{"message": msg})
		return
	}
	if row := config.DB.Save(&desa).RowsAffected; row == 0 {
		ctx.JSON(http.StatusBadRequest, gin.H{"message": "data desa udated failed!"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "data desa udated success!", "desa": &desa})
}
func Delete(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		ctx.Abort()
		return
	}
	var desa Mdesa.Desa
	id := ctx.Param("id")
	check_id := config.DB.Where("id = ?", id).First(&desa)

	if check_id.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Data desa dengan id " + id + " tidak ada !"})
		return
	}
	if err := config.DB.Delete(&desa).Error; err != nil {
		// Error while deleting desa
		ctx.JSON(500, gin.H{"error": "Failed to delete desa"})
		return
	}
	ctx.JSON(200, gin.H{"message": "desa deleted successfully"})
}
