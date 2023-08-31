package dynamodb

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var testTableBasics *TableBasics

// CreateLocalClient Creates a local DynamoDb Client on the specified port. Useful for connecting to DynamoDB Local or
// LocalStack.
func CreateLocalClient(port int) *dynamodb.Client {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("ap-northeast-1"),
		config.WithCredentialsProvider(credentials.StaticCredentialsProvider{
			Value: aws.Credentials{
				AccessKeyID: "dummy", SecretAccessKey: "dummy", SessionToken: "dummy",
				Source: "Hard-coded credentials; values are irrelevant for local DynamoDB",
			},
		}),
	)
	if err != nil {
		panic(err)
	}
	dsn := fmt.Sprintf("http://localhost:%d/", port)
	return dynamodb.NewFromConfig(cfg, func(o *dynamodb.Options) {
		o.BaseEndpoint = aws.String(dsn)
	})
}

func TestMain(m *testing.M) {
	// google_oauth
	testTableBasics = NewTableBasics("google-oauth") // go-dynamodb-reference-table
	// change to local dynamodb
	testTableBasics.DynamoDbClient = CreateLocalClient(8000)
	os.Exit(m.Run())
}

func TestTableExists(t *testing.T) {
	exits, err := testTableBasics.TableExists()
	t.Log("Exists:", exits)
	t.Log("ERROR:", err)
}

func TestListTables(t *testing.T) {
	tableNames, err := testTableBasics.ListTables()
	t.Log("tableNames:", tableNames)
	t.Log("ERROR:", err)

}
