package library

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mirzha99/go-penduduk/helper"
	"github.com/mirzha99/timegoza/timegoza"
)

func UploadFiles(ctx *gin.Context, namefile string, path string) (int, any) {
	file, header, err := ctx.Request.FormFile(namefile)
	if err != nil {
		e := fmt.Sprintf("Error while uploading: %s", err.Error())
		return http.StatusBadRequest, e
	}
	defer file.Close()

	// Pastikan untuk variabel path dengan direktori tempat Anda ingin menyimpan file.
	// Di sini, kita menyimpan file di direktori yang sama dengan aplikasi.
	//replace all space
	namefiles := header.Filename
	namefiles = strings.ReplaceAll(namefiles, " ", "_")
	//name file
	filePath := path + helper.Itoa(int(timegoza.EpochTime())) + "-" + namefiles
	out, err := os.Create(filePath)
	if err != nil {
		e := fmt.Sprintf("Error while creating the file: %s", err.Error())
		return http.StatusInternalServerError, e
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		e := fmt.Sprintf("Error while copying the file: %s", err.Error())
		return http.StatusInternalServerError, e
	}

	//return name file
	return http.StatusCreated, helper.Itoa(int(timegoza.EpochTime())) + "-" + namefiles
}
func RemoveFile(namefile string, path string) {
	filePath := fmt.Sprintf("%v/%v", path, namefile)
	// Hapus file
	err := os.Remove(filePath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
}
