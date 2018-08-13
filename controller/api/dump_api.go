package api

import (
	"github.com/gin-gonic/gin"
	"fmt"
	"github.com/apaxa-go/helper/strconvh"
	"encoding/base64"
	"github.com/sipt/shuttle"
)

func SetAllowDump(ctx *gin.Context) {
	var response Response
	allow_mitm := ctx.PostForm("allow_mitm")
	allow_dump := ctx.PostForm("allow_dump")
	fmt.Println("allow_dump: ", allow_dump)
	switch allow_dump {
	case "true":
		shuttle.SetAllowDump(true)
	case "false":
		shuttle.SetAllowDump(false)
	case "":
	default:
		response.Code = 1
		response.Message = fmt.Sprintf("allow_dump value error: %v", allow_dump)
		ctx.JSON(500, response)
		return
	}
	switch allow_mitm {
	case "true":
		shuttle.SetAllowMitm(true)
	case "false":
		shuttle.SetAllowMitm(false)
	case "":
	default:
		response.Code = 1
		response.Message = fmt.Sprintf("allow_mitm value error: %v", allow_mitm)
		ctx.JSON(500, response)
		return
	}
	GetAllowDump(ctx)
}

func GetAllowDump(ctx *gin.Context) {
	var response = Response{
		Data: struct {
			AllowDump bool `json:"allow_dump"`
			AllowMitm bool `json:"allow_mitm"`
		}{
			shuttle.GetAllowDump(),
			shuttle.GetAllowMitm(),
		},
	}
	ctx.JSON(200, response)
}
func DumpRequest(ctx *gin.Context) {
	var response Response
	idStr := ctx.Param("conn_id")
	id, err := strconvh.ParseInt64(idStr)
	if err != nil {
		response.Code = 1
		response.Message = err.Error()
		ctx.JSON(500, response)
		return
	}
	r := shuttle.GetRecord(id)
	if r == nil {
		response.Code = 1
		response.Message = idStr + " not exist"
		ctx.JSON(500, response)
		return
	}
	if r.Status != shuttle.RecordStatusCompleted {
		response.Code = 1
		response.Message = idStr + " not Completed"
		ctx.JSON(500, response)
		return
	}
	reply := &struct {
		RequestData  string
		ResponseData string
	}{}
	dump := shuttle.GetDump()
	if dump == nil {
		response.Code = 1
		response.Message = "IDump is nil"
		ctx.JSON(500, response)
		return
	}
	data, err := dump.ReadRequest(id)
	if err != nil {
		response.Code = 1
		response.Message = err.Error()
		ctx.JSON(500, response)
		return
	}
	reply.RequestData = base64.RawStdEncoding.EncodeToString(data)
	data, err = dump.ReadResponse(id)
	if err != nil {
		response.Code = 1
		response.Message = err.Error()
		ctx.JSON(500, response)
		return
	}
	reply.ResponseData = base64.RawStdEncoding.EncodeToString(data)
	response.Data = reply
	ctx.JSON(200, response)
}
