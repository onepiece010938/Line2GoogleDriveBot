package app

import (
	"context"

	"github.com/line/line-bot-sdk-go/v7/linebot"

	"github.com/onepiece010938/Line2GoogleDriveBot/internal/adapter/dynamodb"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/adapter/google"
	serviceDrive "github.com/onepiece010938/Line2GoogleDriveBot/internal/app/service/drive"
	serviceSample "github.com/onepiece010938/Line2GoogleDriveBot/internal/app/service/sample"
)

type Application struct {
	SampleService *serviceSample.SampleService
	DriveService  *serviceDrive.GoogleDriveService
	LineBotClient *linebot.Client
}

func NewApplication(ctx context.Context, dynamodb dynamodb.DynamodbI, oauth google.GoogleOAuthI, lineBotClient *linebot.Client) *Application {

	app := &Application{
		LineBotClient: lineBotClient,
		SampleService: serviceSample.NewSampleService(ctx, serviceSample.SampleServiceParam{
			SampleServiceDynamodb: dynamodb,
		}),
		DriveService: serviceDrive.NewGoogleDriveService(ctx, serviceDrive.GoogleDriveServiceParam{
			DriveServiceGoogleOA: oauth,
			DriveServiceDynamodb: dynamodb,
		}),
	}
	return app
}
