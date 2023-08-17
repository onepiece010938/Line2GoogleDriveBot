package app

import (
	"context"

	"github.com/line/line-bot-sdk-go/linebot"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/adapter/cache"
	serviceAnalyze "github.com/onepiece010938/Line2GoogleDriveBot/internal/app/service/analyze"
	serviceMessage "github.com/onepiece010938/Line2GoogleDriveBot/internal/app/service/message"
)

type Application struct {
	// JobService   *serviceJob.JobService
	// ImageService *serviceImage.ImageService
	AnalyzeService *serviceAnalyze.AnalyzeService
	MessageService *serviceMessage.MessageService
	LineBotClient  *linebot.Client
}

func NewApplication(ctx context.Context, cache cache.CacheI, lineBotClient *linebot.Client) *Application {

	// Create application
	app := &Application{
		LineBotClient: lineBotClient,
		MessageService: serviceMessage.NewMessageService(ctx, serviceMessage.MessageServiceParam{
			MessageServiceCache: cache,
		}),
		AnalyzeService: serviceAnalyze.NewAnalyzeService(ctx, serviceAnalyze.AnalyzeServiceParam{
			AnalyzeServiceCache: cache,
		}),
	}
	return app
}
