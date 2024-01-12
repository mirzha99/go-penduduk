package datakk

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/mirzha99/go-penduduk/config"
	"github.com/mirzha99/go-penduduk/helper"
	"github.com/mirzha99/go-penduduk/library"
	Mkepalakularga "github.com/mirzha99/go-penduduk/models/MKepalaKularga"
	"github.com/mirzha99/timegoza/timegoza"
)

// @Summary Get All Data Kepala Keluarga
// @Tags Kepala Keluarga
// @Description Get All Data Kepala Keluarga
// @Security ApiKeyAuth
// @Accept json
// @Router /kepalakeluarga [get]
func Index(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		ctx.Abort()
		return
	}
	var datakk []Mkepalakularga.KepalaKeluargaResult
	config.DB.Table("kepala_keluargas").
		Select("kepala_keluargas.id,kepala_keluargas.nama,kepala_keluargas.nik,kepala_keluargas.id_desa,desas.nama AS nama_desa, desas.nama_kepala_desa AS nama_kepala_desa, desas.id_mukim AS id_mukim,mukims.nama AS nama_mukim, mukims.nama_imum_mukim AS nama_imum_mukim,kepala_keluargas.created_at,kepala_keluargas.change_at,kepala_keluargas.gambar").
		Joins("JOIN desas ON desas.id = kepala_keluargas.id_desa").
		Joins("JOIN mukims ON mukims.id = desas.id_mukim").
		Find(&datakk)

	url := "http://" + ctx.Request.Host
	for i, data := range datakk {
		created, _ := strconv.Atoi(data.Created_at)
		change, _ := strconv.Atoi(data.Change_at)
		convcreated := timegoza.ZaTimes{Epoch: int64(created), TimeZone: "Asia/Jakarta"}
		convchange := timegoza.ZaTimes{Epoch: int64(change), TimeZone: "Asia/Jakarta"}
		datakk[i].Created_at = convcreated.HumanTime()
		datakk[i].Change_at = convchange.HumanTime()
		datakk[i].Gambar = url + "/datapenduduk/" + datakk[i].Gambar
	}
	if len(datakk) == 0 {
		e := helper.ErrorResponse{Detail: "Data Kepala Keluarga Kosong"}
		ctx.JSON(404, e.ErrorResultDetail())
		return
	}
	TokenString, err := ctx.Cookie("jwt-token")
	if err != nil {
		e := helper.ErrorResponse{Detail: "Token Not Found", Error: err.Error(), StatusCode: 401}
		ctx.JSON(401, e.ErrorResultDetail())
		return
	}
	s := helper.SuccessResponse{Result: datakk, StatusCode: 200, Token: TokenString}
	ctx.JSON(200, s.SuccesResult())
}

// @Summary Get All Data Kepala Keluarga
// @Tags Kepala Keluarga
// @Description Get All Data Kepala Keluarga
// @Security ApiKeyAuth
// @Accept json
// @Param id path int true "Kelapa Keluarga id"
// @Router /kepalakeluarga/{id} [get]
func ById(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		ctx.Abort()
		return
	}
	id := ctx.Param("id")
	var datakk Mkepalakularga.KepalaKeluargaResult
	row := config.DB.Table("kepala_keluargas").
		Select("kepala_keluargas.id,kepala_keluargas.nama,kepala_keluargas.nik,kepala_keluargas.id_desa,desas.nama AS nama_desa, desas.nama_kepala_desa AS nama_kepala_desa, desas.id_mukim AS id_mukim,mukims.nama AS nama_mukim, mukims.nama_imum_mukim AS nama_imum_mukim,kepala_keluargas.created_at,kepala_keluargas.change_at,kepala_keluargas.gambar").
		Joins("JOIN desas ON desas.id = kepala_keluargas.id_desa").
		Joins("JOIN mukims ON mukims.id = desas.id_mukim").
		Where("kepala_keluargas.id = ?", id).
		Find(&datakk)

	created, _ := strconv.Atoi(datakk.Created_at)
	change, _ := strconv.Atoi(datakk.Change_at)
	convcreated := timegoza.ZaTimes{Epoch: int64(created), TimeZone: "Asia/Jakarta"}
	convchange := timegoza.ZaTimes{Epoch: int64(change), TimeZone: "Asia/Jakarta"}
	datakk.Created_at = convcreated.HumanTime()
	datakk.Change_at = convchange.HumanTime()
	if row.RowsAffected == 0 {
		e := helper.ErrorResponse{Detail: "Data Kepala Keluarga Kosong", StatusCode: 404}
		ctx.JSON(404, e.ErrorResultDetail())
		return
	}
	TokenString, err := ctx.Cookie("jwt-token")
	if err != nil {
		e := helper.ErrorResponse{Detail: "Token Not Found", Error: err.Error(), StatusCode: 401}
		ctx.JSON(401, e.ErrorResultDetail())
		return
	}
	s := helper.SuccessResponse{Result: datakk, StatusCode: 200, Token: TokenString}
	ctx.JSON(200, s.SuccesResult())
}
func Add(ctx *gin.Context) {
	//deklarasi form post
	nama := ctx.PostForm("nama")
	nik := helper.Atoi(ctx.PostForm("nik")).(int)
	id_desa := helper.Atoi(ctx.PostForm("nik")).(int)
	//struct kepala keluarga berdasarkan data form post
	datakk := Mkepalakularga.KepalaKeluarga{
		Nama:    nama,
		Nik:     nik,
		Id_Desa: id_desa,
	}
	//ubah data struct
	datakk.Created_at = helper.Itoa(int(timegoza.EpochTime()))
	datakk.Change_at = helper.Itoa(int(timegoza.EpochTime()))
	//check ShouldBind
	if err := ctx.ShouldBind(&datakk); err != nil {
		e := helper.ErrorResponse{Detail: err.Error(), Error: "Should Bind", StatusCode: http.StatusBadRequest}
		ctx.JSON(http.StatusBadRequest, e.ErrorResultDetail())
		return
	}
	statuscode, result := library.UploadFiles(ctx, "gambar", "uploads/datapenduduk/")
	if statuscode != http.StatusCreated {
		e := helper.ErrorResponse{Detail: result.(string), Error: "Error Upload Image", StatusCode: statuscode}
		ctx.JSON(http.StatusBadRequest, e.ErrorResultDetail())
		return
	}
	datakk.Gambar = result.(string)
	row := config.DB.Create(&datakk)
	if row.RowsAffected == 0 {
		e := helper.ErrorResponse{Detail: "Data Kepala Keluarga gagal ditambah", Error: "Bad Request", StatusCode: http.StatusBadRequest}
		ctx.JSON(http.StatusBadRequest, e.ErrorResultDetail())
		return
	}
	s := helper.SuccessResponse{Message: "Data Kepala Keluarga Berhasil Di Tambah", Result: datakk, StatusCode: http.StatusCreated}
	ctx.JSON(s.StatusCode, s.SuccesResult())
}
