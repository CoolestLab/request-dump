package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/CoolestLab/request-dump/build"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	SLACK_API = os.Getenv("SLACK_API")
)

func Execute() {
	if err := server(); err != nil {
		panic(err)
	}
}

func server() error {
	e := gin.Default()

	e.Use(cors.Default())
	e.NoRoute(errFuncWrapper(func(c *gin.Context) (interface{}, error) {
		return c.Request.URL.Path, errors.New("not found")
	}))

	e.GET("/version", errFuncWrapper(func(c *gin.Context) (interface{}, error) {
		return build.InfoMap, nil
	}))
	e.Any("/dump", errFuncWrapper(func(c *gin.Context) (interface{}, error) {
		if len(SLACK_API) == 0 {
			return nil, errors.New("no slack api configured")
		}

		payload := new(strings.Builder)
		payload.WriteString(fmt.Sprintf("method: %s\n", c.Request.Method))
		payload.WriteString(fmt.Sprintf("uri: %s\n", c.Request.URL.String()))
		payload.WriteString("headers:\n")
		for k, v := range c.Request.Header {
			payload.WriteString(fmt.Sprintf("  %s: %s\n", k, strings.Join(v, ",")))
		}
		payload.WriteString("body:\n")
		body, _ := io.ReadAll(c.Request.Body)
		payload.WriteString(string(body) + "\n")

		data, _ := json.Marshal(struct {
			Text string `json:"text"`
		}{payload.String()})

		resp, err := http.Post(SLACK_API, "application/json", bytes.NewBuffer(data))
		if err != nil {
			return nil, err
		}

		return resp.StatusCode, nil
	}))

	gin.SetMode(gin.ReleaseMode)
	return e.Run(":9000")
}

func errFuncWrapper(f func(*gin.Context) (interface{}, error)) func(*gin.Context) {
	return func(ctx *gin.Context) {
		resp, err := f(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"code": 1,
				"msg":  err.Error(),
				"data": resp,
			})
			return
		}

		ctx.AbortWithStatusJSON(http.StatusOK, gin.H{
			"code": 0,
			"msg":  "success",
			"data": resp,
		})
	}
}
