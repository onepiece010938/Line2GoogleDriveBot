package httputil

import (
	"github.com/gin-gonic/gin"
)

// NewError example
func NewError(ctx *gin.Context, status int, err_code int, append_msg *string) {
	if append_msg == nil {
		empty := ""
		append_msg = &empty
	} else {
		*append_msg = " <append_msg>:" + *append_msg
	}
	er := HTTPError{
		Status:  status,
		Code:    err_code,
		Message: *append_msg,
	}
	ctx.JSON(status, er)
}

// HTTPError example
type HTTPError struct {
	Status  int    `json:"status" example:"400"`
	Code    int    `json:"code" example:"40012"`
	Message string `json:"msg" example:"status bad request"`
}
