package sample

import "github.com/onepiece010938/Line2GoogleDriveBot/internal/adapter/dynamodb"

type SampleServiceDynamodbI interface {
	GetGoogleOAuthToken(line_userid string) (dynamodb.GoogleOAuthToken, error)
}
