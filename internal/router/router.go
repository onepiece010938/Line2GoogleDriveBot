package router

import (
	"context"

	"github.com/onepiece010938/Line2GoogleDriveBot/internal/app"
	v1 "github.com/onepiece010938/Line2GoogleDriveBot/internal/router/api/v1"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/router/middleware"

	"github.com/gin-gonic/gin"
)

/*
func InitRouter(ctx context.Context, app *app.Application) *gin.Engine {
	// docs.SwaggerInfo.BasePath = "/api/v1"

	if viper.GetBool("release") {
		gin.SetMode(gin.ReleaseMode)
	}
	r := gin.New()
	// r.Use(gin.Logger())
	// r.Use(gin.Recovery())

	// r.Use(log.GinLogger(log.Logger), log.GinRecovery(log.Logger, true))

	//use Cors
	r.Use(middleware.Cors())

	// use swagger (close in ReleaseMode)
	// if gin.Mode() != gin.ReleaseMode {
	// 	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	// 	// http://127.0.0.1:7777/swagger/index.html
	// }

	// for handler Response
	// responser := NewResponder(ctx)

	// gin.SetMode(setting.RunMode)
	setUpRouter(r, app)

	return r
}

func setUpRouter(router *gin.Engine, app *app.Application) {
	api := router.Group("/api")
	{
		v1.RegisterRouter(api, app)
		//v2.RegisterRouter(api)
	}
}
*/

func SetGeneralMiddlewares(ctx context.Context, ginRouter *gin.Engine) {
	ginRouter.Use(gin.Recovery())
	ginRouter.Use(middleware.Cors())

}
func RegisterHandlers(router *gin.Engine, app *app.Application) {
	registerAPIHandlers(router, app)
}

func registerAPIHandlers(router *gin.Engine, app *app.Application) {
	// Build middlewares
	// BearerToken := NewAuthMiddlewareBearer(app)

	api := router.Group("/api")
	{
		v1.RegisterRouter(api, app)
		//v2.RegisterRouter(api)
	}

}
