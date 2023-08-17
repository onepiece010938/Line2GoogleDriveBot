package v1

import (
	"fmt"
	"net/http"

	"github.com/onepiece010938/Line2GoogleDriveBot/internal/app"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/app/service/analyze"

	"github.com/gin-gonic/gin"
)

func StartAnalyze(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		err := app.AnalyzeService.CreateAnalyze(ctx, analyze.CreateAnalyzeParm{})
		if err != nil {
			fmt.Println(err)
		}
		c.JSON(http.StatusOK, "")
	}

}
