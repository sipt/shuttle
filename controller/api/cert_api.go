package api

import (
	"github.com/gin-gonic/gin"
	"io"
	"bytes"
	"github.com/sipt/shuttle"
)

func DownloadCert(ctx *gin.Context) {
	var response Response
	caBytes := shuttle.GetCACert()
	if len(caBytes) == 0 {
		response.Code = 1
		response.Message = "please generate CA"
		ctx.JSON(500, response)
		return
	}
	ctx.Header("Content-Type", "application/octet-stream")
	ctx.Header("content-disposition", "attachment; filename=\"Shuttle.cer\"")
	_, err := io.Copy(ctx.Writer, bytes.NewBuffer(caBytes))
	if err != nil {
		response.Code = 1
		response.Message = err.Error()
		ctx.JSON(500, response)
		return
	}
}
func GenerateCert(ctx *gin.Context) {
	var response Response
	err := shuttle.GenerateCA()
	if err != nil {
		response.Code = 1
		response.Message = err.Error()
		ctx.JSON(500, response)
		return
	}
	ctx.JSON(200, response)
}
