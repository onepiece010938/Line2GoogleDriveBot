package drive

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	dynamodbAdapter "github.com/onepiece010938/Line2GoogleDriveBot/internal/adapter/dynamodb"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/adapter/google"
	"golang.org/x/oauth2"
)

type DriveServiceGoogleOAuthI interface {
	OAuthLoginURL(lineID string) (oauthURL string)
	UserOAuthToken(authCode string) (*oauth2.Token, error)
	NewGoogleDrive(ctx context.Context, tok *oauth2.Token) (*google.GoogleDrive, error)
}
type DriveServiceDynamodbI interface {
	GetGoogleOAuthToken(line_userid string) (dynamodbAdapter.GoogleOAuthToken, error)
	CreateGoogleOAuthTable() (*types.TableDescription, error)
	AddGoogleOAuthToken(tok dynamodbAdapter.GoogleOAuthToken) error
	TxUpdateGoogleOAuthToken(tok dynamodbAdapter.GoogleOAuthToken) (*dynamodb.TransactWriteItemsOutput, error)
}
