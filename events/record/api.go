package record

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/controller/model"
)

func InitAPI(e *gin.Engine) {
	e.GET("/api/records", recordsHandleFunc)
}

func recordsHandleFunc(c *gin.Context) {
	list := make([]RecordEntity, 0, recordStarge.Len())
	recordStarge.Range(func(v interface{}) bool {
		r := v.(*RecordEntity)
		list = append(list, *r)
		return false
	})
	c.JSON(http.StatusOK, &model.Response{
		Code: 0,
		Data: list,
	})
}
