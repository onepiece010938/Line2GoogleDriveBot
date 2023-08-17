package v1

import (
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/app"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(router *gin.RouterGroup, app *app.Application) {
	v1 := router.Group("/v1")
	{
		v1.POST("/callback", Callback(app))
		v1.GET("/sample", SAMPLE)
		v1.GET("/analyze", StartAnalyze(app))
	}
}
