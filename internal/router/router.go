package router

import (
	"context"

	"github.com/onepiece010938/Line2GoogleDriveBot/internal/app"
	v1 "github.com/onepiece010938/Line2GoogleDriveBot/internal/router/api/v1"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/router/middleware"

	"github.com/gin-gonic/gin"
)

func SetGeneralMiddlewares(ctx context.Context, ginRouter *gin.Engine) {
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(middleware.Cors())

}
func RegisterHandlers(router *gin.Engine, app *app.Application) {
	registerAPIHandlers(router, app)
}

func registerAPIHandlers(router *gin.Engine, app *app.Application) {

	api := router.Group("/api")
	{
		v1.RegisterRouter(api, app)
		//v2.RegisterRouter(api)
	}

}
