package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SAMPLE(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, "SAMPLE")
}
