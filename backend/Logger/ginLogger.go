package logger

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func abortErr(c *gin.Context, code int, msg string, err error) {
	LogErr("%v", err)
	c.AbortWithStatusJSON(code, msg)
}

func abort(c *gin.Context, code int, msg string) {
	LogErr("%s", msg)
	c.AbortWithStatusJSON(code, msg)
}

func BadRequestStr(c *gin.Context, msg string) {
	abort(c, http.StatusBadRequest, msg)
}

func InternalServerErrorStr(c *gin.Context, msg string) {
	abort(c, http.StatusInternalServerError, msg)
}

func InternalServerError(c *gin.Context, err error) {
	abortErr(c, http.StatusInternalServerError, err.Error(), err)
}
