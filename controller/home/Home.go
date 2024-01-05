package home

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mirzha99/go-penduduk/config"
)

func Index(ctx *gin.Context) {
	if !config.Limiter.Allow() {
		ctx.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests"})
		ctx.Abort()
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"Message": "Hello World !"})
}
