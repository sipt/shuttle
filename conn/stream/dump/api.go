package dump

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"os"
	"path"
	"strconv"

	"github.com/sipt/shuttle/events/record"

	"github.com/gin-gonic/gin"
	"github.com/sipt/shuttle/constant/typ"
	"github.com/sipt/shuttle/controller/model"
	"github.com/sipt/shuttle/global/namespace"
	"github.com/sirupsen/logrus"
)

func InitAPI(e *gin.Engine) {
	r := e.Group("/api/stream/dump")
	r.GET("/status", getStatus)
	r.PUT("/status", putStatus)
	r.POST("/cert", generateCA)
	r.GET("/cert", downloadCA)

	r.GET("/session/:id", dumpSession)
	r.GET("/request/body/:id", dumpRequest)
	r.GET("/response/body/:id", dumpResponse)
}

func getStatus(c *gin.Context) {
	c.JSON(200, &model.Response{
		Data: gin.H{
			"dump": allowDump,
			"mitm": mitmEnabled,
		},
	})
}

func putStatus(c *gin.Context) {
	data, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(500, &model.Response{
			Code:    1,
			Message: err.Error(),
		})
		return
	}
	entity := make(map[string]bool)
	err = json.Unmarshal(data, &entity)
	if err != nil {
		c.JSON(400, &model.Response{
			Code:    1,
			Message: err.Error(),
		})
		return
	}
	runtime := namespace.NamespaceWithName("default").Runtime()
	runtime = typ.NewRuntime("stream", runtime)
	runtime = typ.NewRuntime("data-dump", runtime)
	if v, ok := entity["dump"]; ok && allowDump != v {
		allowDump = v
		err = runtime.Set("dump", v)
		if err != nil {
			logrus.WithError(err).Error("set stream.data-dump.dump failed")
		}
	}
	if v, ok := entity["mitm"]; ok && mitmEnabled != v {
		mitmEnabled = v
		err = runtime.Set("mitm", v)
		if err != nil {
			logrus.WithError(err).Error("set stream.data-dump.mitm failed")
		}
	}
	c.JSON(200, &model.Response{})
}

func generateCA(c *gin.Context) {
	key, ca, err := GenerateCA()
	if err != nil {
		c.JSON(500, &model.Response{
			Code:    1,
			Message: err.Error(),
		})
		return
	}
	runtime := namespace.NamespaceWithName("default").Runtime()
	runtime = typ.NewRuntime("stream", runtime)
	runtime = typ.NewRuntime("data-dump", runtime)
	err = runtime.Set("key", key)
	if err != nil {
		logrus.WithError(err).Error("set stream.data-dump.key failed")
	}
	err = runtime.Set("ca", ca)
	if err != nil {
		logrus.WithError(err).Error("set stream.data-dump.ca failed")
	}
	c.JSON(200, &model.Response{})
}

func downloadCA(c *gin.Context) {
	if len(caBytes) == 0 {
		c.JSON(500, &model.Response{
			Code:    1,
			Message: "please generate CA first",
		})
		return
	}
	bak := make([]byte, len(caBytes))
	copy(bak, caBytes)
	c.Header("Content-Type", "application/octet-stream")
	c.Header("content-disposition", "attachment; filename=\"Shuttle.cer\"")
	_, err := io.Copy(c.Writer, bytes.NewBuffer(bak))
	if err != nil {
		c.JSON(500, &model.Response{
			Code:    1,
			Message: err.Error(),
		})
		return
	}
	c.JSON(200, &model.Response{})
}

func dumpSession(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(400, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("[%s] not a number", idStr),
		})
		return
	}
	found, status, dumped := false, record.ActiveStatus, false
	record.RangeRecord(c, func(entity *record.RecordEntity) bool {
		if entity.ID == id {
			found = true
			status = entity.Status
			dumped = entity.Dumped
			return true
		}
		return false
	})
	if !found {
		c.JSON(400, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("invalid recordID [%s]", idStr),
		})
		return
	}
	if status != record.CompletedStatus {
		c.JSON(400, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("record[%s] is not complete", idStr),
		})
		return
	}
	reqFile, err := os.Open(path.Join(dirPath, ReqKey(id)))
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(400, &model.Response{
				Code:    1,
				Message: fmt.Sprintf("record[%s] is not dumped", idStr),
			})
			return
		} else {
			c.JSON(500, &model.Response{
				Code:    1,
				Message: fmt.Sprintf("record[%s] get request data failed", idStr),
			})
			return
		}
	}
	defer reqFile.Close()
	respFile, err := os.Open(path.Join(dirPath, RespKey(id)))
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(400, &model.Response{
				Code:    1,
				Message: fmt.Sprintf("record[%s] is not dumped", idStr),
			})
			return
		} else {
			c.JSON(500, &model.Response{
				Code:    1,
				Message: fmt.Sprintf("record[%s] get response data failed", idStr),
			})
			return
		}
	}
	defer respFile.Close()
	req, err := http.ReadRequest(bufio.NewReader(reqFile))
	if err != nil {
		c.JSON(500, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("record[%s] get request data failed", idStr),
		})
		return
	}
	resp, err := http.ReadResponse(bufio.NewReader(respFile), req)
	if err != nil {
		c.JSON(500, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("record[%s] get response data failed", idStr),
		})
		return
	}
	w := base64.NewEncoder(base64.StdEncoding, c.Writer)
	_, _ = c.Writer.Write([]byte(`{"request":{"header":"`))
	b, _ := httputil.DumpRequest(req, false)
	_, _ = w.Write(b)
	_ = w.Close()
	_, _ = c.Writer.Write([]byte(`","body":"`))
	_, _ = io.Copy(w, req.Body)
	_ = w.Close()
	_, _ = c.Writer.Write([]byte(`"}`))
	// write response
	_, _ = c.Writer.Write([]byte(`, "response":{"header":"`))
	b, _ = httputil.DumpResponse(resp, false)
	_, _ = w.Write(b)
	_ = w.Close()
	_, _ = c.Writer.Write([]byte(`","body":"`))
	_, _ = io.Copy(w, resp.Body)
	_ = w.Close()
	_, _ = c.Writer.Write([]byte(`"}`))
	c.Writer.WriteHeader(200)
}
func dumpRequest(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(400, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("[%s] not a number", idStr),
		})
		return
	}
	found, status, dumped := false, record.ActiveStatus, false
	record.RangeRecord(c, func(entity *record.RecordEntity) bool {
		if entity.ID == id {
			found = true
			status = entity.Status
			dumped = entity.Dumped
			return true
		}
		return false
	})
	if !found {
		c.JSON(400, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("invalid recordID [%s]", idStr),
		})
		return
	}
	if status != record.CompletedStatus {
		c.JSON(400, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("record[%s] is not complete", idStr),
		})
		return
	}
	reqFile, err := os.Open(path.Join(dirPath, ReqKey(id)))
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(400, &model.Response{
				Code:    1,
				Message: fmt.Sprintf("record[%s] is not dumped", idStr),
			})
			return
		} else {
			c.JSON(500, &model.Response{
				Code:    1,
				Message: fmt.Sprintf("record[%s] get request data failed", idStr),
			})
			return
		}
	}
	defer reqFile.Close()
	req, err := http.ReadRequest(bufio.NewReader(reqFile))
	if err != nil {
		c.JSON(500, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("record[%s] get request data failed", idStr),
		})
		return
	}
	io.Copy(c.Writer, req.Body)
	c.Writer.WriteHeader(200)
}
func dumpResponse(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(400, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("[%s] not a number", idStr),
		})
		return
	}
	found, status, dumped := false, record.ActiveStatus, false
	record.RangeRecord(c, func(entity *record.RecordEntity) bool {
		if entity.ID == id {
			found = true
			status = entity.Status
			dumped = entity.Dumped
			return true
		}
		return false
	})
	if !found {
		c.JSON(400, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("invalid recordID [%s]", idStr),
		})
		return
	}
	if status != record.CompletedStatus {
		c.JSON(400, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("record[%s] is not complete", idStr),
		})
		return
	}
	respFile, err := os.Open(path.Join(dirPath, RespKey(id)))
	if err != nil {
		if os.IsNotExist(err) {
			c.JSON(400, &model.Response{
				Code:    1,
				Message: fmt.Sprintf("record[%s] is not dumped", idStr),
			})
			return
		} else {
			c.JSON(500, &model.Response{
				Code:    1,
				Message: fmt.Sprintf("record[%s] get response data failed", idStr),
			})
			return
		}
	}
	defer respFile.Close()
	resp, err := http.ReadResponse(bufio.NewReader(respFile), &http.Request{})
	if err != nil {
		c.JSON(500, &model.Response{
			Code:    1,
			Message: fmt.Sprintf("record[%s] get response data failed", idStr),
		})
		return
	}
	io.Copy(c.Writer, resp.Body)
	c.Writer.WriteHeader(200)
}
