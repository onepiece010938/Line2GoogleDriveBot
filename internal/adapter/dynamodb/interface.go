package dynamodb

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type DynamodbI interface {
	CreateGoogleOAuthTable() (*types.TableDescription, error)
	AddGoogleOAuthToken(tok GoogleOAuthToken) error
	TxUpdateGoogleOAuthToken(tok GoogleOAuthToken) (*dynamodb.TransactWriteItemsOutput, error)
	GetGoogleOAuthToken(line_userid string) (GoogleOAuthToken, error)
}
