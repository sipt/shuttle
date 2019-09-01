package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

const (
	ParamsKeyAuthType = "auth_type"

	AuthTypeBasic  = "basic"
	AuthTypeBearer = "bearer"

	ParamsKeyUser     = "user"
	ParamsKeyPassword = "password"
	ParamsKeyToken    = "token"
)

var e *gin.Engine

func InitEngine(addr string, params map[string]string) (closer func(), err error) {
	e = gin.Default()
	authType, ok := params[ParamsKeyAuthType]
	if ok {
		var authFunc gin.HandlerFunc
		switch authType {
		case AuthTypeBasic:
			authFunc, err = newBasicAuth(params)
		case AuthTypeBearer:
			authFunc, err = newBearerAuth(params)
		default:
			err = errors.Errorf("[http.Inbound] is not support")
		}
		if err != nil {
			return
		}
		e.Use(authFunc)
	}
	server := &http.Server{
		Addr:           addr,
		Handler:        e,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go server.ListenAndServe()
	return func() {
		ctx, _ := context.WithTimeout(context.Background(), 3*time.Second)
		server.Shutdown(ctx)
	}, nil
}

func newBasicAuth(params map[string]string) (gin.HandlerFunc, error) {
	user := params[ParamsKeyUser]
	if len(user) == 0 {
		return nil, errors.New("[user] is empty")
	}
	password := params[ParamsKeyPassword]
	if len(password) == 0 {
		return nil, errors.New("[password] is empty")
	}
	return gin.BasicAuth(gin.Accounts{user: password}), nil
}

func newBearerAuth(params map[string]string) (gin.HandlerFunc, error) {
	token := params[ParamsKeyToken]
	if len(token) == 0 {
		return nil, errors.New("[token] is empty")
	}
	authorization := "Bearer " + token
	return func(c *gin.Context) {
		if c.GetHeader("Authorization") != authorization {
			c.JSON(http.StatusUnauthorized, &Response{
				Code:    1,
				Message: "unauthorized",
			})
		}
	}, nil
}

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
