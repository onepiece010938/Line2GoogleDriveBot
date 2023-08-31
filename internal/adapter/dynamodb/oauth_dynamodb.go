package dynamodb

import (
	"context"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type GoogleOAuthToken struct {
	PK           string `dynamodbav:"PK"`
	AccessToken  string `dynamodbav:"access_token"`
	TokenType    string `dynamodbav:"token_type"`
	RefreshToken string `dynamodbav:"refresh_token"`
	Expiry       string `dynamodbav:"expiry"`

	Info map[string]interface{} `dynamodbav:"info"`
}

const PK_PREFIX_LINE = "LINEID#"

// Get PK Key (with Line prefix)
func (tok GoogleOAuthToken) GetKey() map[string]types.AttributeValue {
	// Add prefix to PK
	line_userid, err := attributevalue.Marshal(PK_PREFIX_LINE + tok.PK)
	if err != nil {
		panic(err)
	}

	return map[string]types.AttributeValue{"PK": line_userid}
}

func (basics TableBasics) CreateGoogleOAuthTable() (*types.TableDescription, error) {
	var tableDesc *types.TableDescription
	table, err := basics.DynamoDbClient.CreateTable(context.TODO(), &dynamodb.CreateTableInput{

		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("PK"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{{
			AttributeName: aws.String("PK"),
			KeyType:       types.KeyTypeHash,
		}},
		TableName: aws.String(basics.TableName),
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(10),
			WriteCapacityUnits: aws.Int64(10),
		},
	})
	if err != nil {
		log.Printf("Couldn't create table %v. Here's why: %v\n", basics.TableName, err)
	} else {
		waiter := dynamodb.NewTableExistsWaiter(basics.DynamoDbClient)
		err = waiter.Wait(context.TODO(), &dynamodb.DescribeTableInput{
			TableName: aws.String(basics.TableName)}, 5*time.Minute)
		if err != nil {
			log.Printf("Wait for table exists failed. Here's why: %v\n", err)
		}
		tableDesc = table.TableDescription
	}
	return tableDesc, err
}

func (basics TableBasics) AddGoogleOAuthToken(tok GoogleOAuthToken) error {
	tok.PK = PK_PREFIX_LINE + tok.PK
	item, err := attributevalue.MarshalMap(tok)
	if err != nil {
		panic(err)
	}
	_, err = basics.DynamoDbClient.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName: aws.String(basics.TableName), Item: item,
	})
	if err != nil {
		log.Printf("Couldn't add item to table. Here's why: %v\n", err)
	}
	return err
}

func (basics TableBasics) UpdateGoogleOAuthToken(tok GoogleOAuthToken) (map[string]string, error) {
	var err error
	var response *dynamodb.UpdateItemOutput
	var attributeMap map[string]string
	update := expression.Set(expression.Name("refresh_token"), expression.Value(tok.RefreshToken))
	update.Set(expression.Name("access_token"), expression.Value(tok.AccessToken))
	// update := expression.Set(expression.Name("info.rating"), expression.Value(tok.Info["rating"]))
	// update.Set(expression.Name("info.plot"), expression.Value(tok.Info["plot"]))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Printf("Couldn't build expression for update. Here's why: %v\n", err)
	} else {
		response, err = basics.DynamoDbClient.UpdateItem(context.TODO(), &dynamodb.UpdateItemInput{
			TableName:                 aws.String(basics.TableName),
			Key:                       tok.GetKey(),
			ExpressionAttributeNames:  expr.Names(),
			ExpressionAttributeValues: expr.Values(),
			UpdateExpression:          expr.Update(),
			ReturnValues:              types.ReturnValueUpdatedNew,
		})
		if err != nil {
			log.Printf("Couldn't update tok %v. Here's why: %v\n", tok.PK, err)
		} else {
			log.Println(response.Attributes)
			err = attributevalue.UnmarshalMap(response.Attributes, &attributeMap)
			if err != nil {
				log.Printf("Couldn't unmarshall update response. Here's why: %v\n", err)
			}
		}
	}
	log.Println(attributeMap)
	return attributeMap, err
}

func (basics TableBasics) TxUpdateGoogleOAuthToken(tok GoogleOAuthToken) (*dynamodb.TransactWriteItemsOutput, error) {
	var err error
	var response *dynamodb.TransactWriteItemsOutput

	update := expression.Set(expression.Name("refresh_token"), expression.Value(tok.RefreshToken))
	update.Set(expression.Name("access_token"), expression.Value(tok.AccessToken))
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		log.Printf("Couldn't build expression for update. Here's why: %v\n", err)
	} else {
		twii := &dynamodb.TransactWriteItemsInput{
			TransactItems: []types.TransactWriteItem{
				{
					Update: &types.Update{
						Key:                       tok.GetKey(),
						TableName:                 aws.String(basics.TableName),
						ExpressionAttributeNames:  expr.Names(),
						ExpressionAttributeValues: expr.Values(),
						UpdateExpression:          expr.Update(),
					},
				},
			},
		}
		response, err = basics.DynamoDbClient.TransactWriteItems(context.TODO(), twii)
		if err != nil {
			log.Printf("Couldn't trasnaciton update tok %v. Here's why: %v\n", tok.PK, err)
		}
	}

	return response, err
}

func (basics TableBasics) GetGoogleOAuthToken(line_userid string) (GoogleOAuthToken, error) {
	tok := GoogleOAuthToken{PK: line_userid}
	response, err := basics.DynamoDbClient.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: tok.GetKey(), TableName: aws.String(basics.TableName),
	})
	if err != nil {
		log.Printf("Couldn't get info about %v. Here's why: %v\n", line_userid, err)
	} else {
		err = attributevalue.UnmarshalMap(response.Item, &tok)
		if err != nil {
			log.Printf("Couldn't unmarshal response. Here's why: %v\n", err)
		}
	}
	return tok, err
}
