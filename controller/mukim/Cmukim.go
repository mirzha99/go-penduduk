package mukim

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mirzha99/go-penduduk/config"
	"github.com/mirzha99/go-penduduk/models/Mdesa"
	"github.com/mirzha99/go-penduduk/models/Mmukim"
)

// @Summary Get All Mukim
// @Tags Mukim
// @Description Get All Mukim
// @Security ApiKeyAuth
// @Accept json
// @Success 201 {object} helper.SuccessResponse
// @Failure 400 {object} helper.ErrorResponse
// @Router /mukims [get]
func Index(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		ctx.Abort()
		return
	}
	var mukim []Mmukim.Mukim
	config.DB.Find(&mukim)

	if len(mukim) == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Mukim is Empty"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"mukim": mukim})
}

// @Summary Get Mukim by Id
// @Tags Mukim
// @Description Get Mukim by Id
// @Security ApiKeyAuth
// @Accept json
// @Param id path int true "User ID"
// @Success 201 {object} helper.SuccessResponse
// @Failure 400 {object} helper.ErrorResponse
// @Router /mukim/{id} [get]
func GetById(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		ctx.Abort()
		return
	}
	var mukim Mmukim.Mukim
	id := ctx.Param("id")
	row := config.DB.Where("id = ?", id).First(&mukim)

	if row.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Data mukim dengan id " + id + " tidak ada !"})
		return
	}

	ctx.JSON(200, gin.H{"mukim": mukim})
}
func mukim_already_exits(nama_mukim string) bool {
	var mukim Mmukim.Mukim
	if err := config.DB.Where("nama = ?", nama_mukim).First(&mukim).Error; err != nil {
		return false
	}
	return true
}

// @Summary Mukim Edit Add
// @Tags Mukim
// @Description Mukim Edit Add
// @Security ApiKeyAuth
// @Accept json
// @Param user body Mmukim.Mukim true "User information"
// @Success 201 {object} helper.SuccessResponse
// @Failure 400 {object} helper.ErrorResponse
// @Router /mukim/ [post]
func Add(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		ctx.Abort()
		return
	}
	var mukim Mmukim.Mukim
	if err := ctx.ShouldBindJSON(&mukim); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if mukim_already_exits(mukim.Nama) {
		ctx.JSON(http.StatusConflict, gin.H{"Message": "Nama Mukim Duplicate!"})
		return
	}
	result := config.DB.Create(&mukim)
	if result.Error != nil {
		ctx.JSON(400, gin.H{"message": "Mukim created Failed", "user": mukim})
		return
	}
	ctx.JSON(200, gin.H{"message": "Mukim created Success", "mukim": mukim})

}

// @Summary Mukim Edit by Id
// @Tags Mukim
// @Description Mukim Edit by Id
// @Security ApiKeyAuth
// @Accept json
// @Param id path int true "User ID"
// @Param user body Mmukim.Mukim true "User information"
// @Success 201 {object} helper.SuccessResponse
// @Failure 400 {object} helper.ErrorResponse
// @Router /mukim/{id} [put]
func Edit(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		ctx.Abort()
		return
	}
	var mukim Mmukim.Mukim
	id := ctx.Param("id")
	check_id := config.DB.Where("id = ?", id).First(&mukim)

	if check_id.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Data mukim dengan id " + id + " tidak ada !"})
		return
	}

	if err := ctx.ShouldBindJSON(&mukim); err != nil {
		ctx.AbortWithStatus(http.StatusBadRequest)
		return
	}

	row := config.DB.Save(&mukim)
	if row.RowsAffected == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "User update failed"})
		return
	}
	ctx.JSON(200, gin.H{"message": "User Successly Update", "user": mukim})
}

// @Summary Mukim Edit by Id
// @Tags Mukim
// @Description Mukim Edit by Id
// @Security ApiKeyAuth
// @Accept json
// @Param id path int true "User ID"
// @Success 201 {object} helper.SuccessResponse
// @Failure 400 {object} helper.ErrorResponse
// @Router /mukim/{id} [delete]
func Delete(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		ctx.Abort()
		return
	}
	var mukim Mmukim.Mukim
	id := ctx.Param("id")
	check_id := config.DB.Where("id = ?", id).First(&mukim)

	if check_id.RowsAffected == 0 {
		ctx.JSON(http.StatusNotFound, gin.H{"message": "Data mukim dengan id " + id + " tidak ada !"})
		return
	}
	if err := config.DB.Delete(&mukim).Error; err != nil {
		// Error while deleting mukim
		ctx.JSON(500, gin.H{"error": "Failed to delete mukim"})
		return
	}
	var desa Mdesa.Desa
	config.DB.Where("id_mukim = ?", id).Delete(&desa)
	ctx.JSON(200, gin.H{"message": "Mukim dan desa bermukim " + mukim.Nama + " berhasil terhapus"})
}
