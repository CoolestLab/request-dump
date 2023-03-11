package cmd

import (
	"errors"
	"net/http"

	"github.com/CoolestLab/serverless-go-template/build"

	"github.com/gin-gonic/contrib/cors"
	"github.com/gin-gonic/gin"
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
