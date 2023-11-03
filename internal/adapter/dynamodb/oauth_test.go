package dynamodb

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCreateGoogleOAuthTable(t *testing.T) {
	tableDesc, err := testTableBasics.CreateGoogleOAuthTable()
	t.Log("tableDesc:", tableDesc)
	t.Log("ERROR:", err)
	assert.NoError(t, err, "Expected no error creating table")
	assert.NotNil(t, tableDesc, "Table description should not be nil")
	// assert.Equal(t, "ExpectedTableName", tableDesc.TableName, "Table name mismatch")
	// if err ==
	// var notFoundEx *types.ResourceNotFoundException
}

func TestAddGoogleOAuthToken(t *testing.T) {

	dateString := "2023-08-22T11:31:54.2936004+08:00"
	// 使用 time.Parse 函數進行轉換
	parsedTime, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		fmt.Println("日期轉換錯誤:", err)
		return
	}

	// 輸出轉換後的 time.Time
	fmt.Println("轉換後的時間:", parsedTime)

	tok := GoogleOAuthToken{
		PK:           "test123",
		AccessToken:  "test123",
		TokenType:    "Bearer",
		RefreshToken: "test123",
		Expiry:       parsedTime,
		Info: map[string]interface{}{
			"upload_folder_id": ""},
	}
	err = testTableBasics.AddGoogleOAuthToken(tok)
	if err != nil {
		t.Log("ERROR:", err)
	}
}

func TestUpdateGoogleOAuthToken(t *testing.T) {
	tok := GoogleOAuthToken{
		PK:          "test1234",
		AccessToken: "test12345",
		// TokenType:    "Bearer",
		RefreshToken: "test12345",
		// Expiry:       "2023-08-22T11:31:54.2936004+08:00",
	}
	output, err := testTableBasics.UpdateGoogleOAuthToken(tok)
	t.Log("output:", output)
	if err != nil {
		t.Log("ERROR:", err)
	}
}

func TestTxUpdateGoogleOAuthToken(t *testing.T) {

	tok := GoogleOAuthToken{
		PK:           "test123",
		AccessToken:  "test123456",
		RefreshToken: "test123456",
		TokenType:    "Bearer2",
		Expiry:       time.Now(),
		Info: map[string]interface{}{
			"upload_folder_id": "GGGGG"},
	}
	output, err := testTableBasics.TxUpdateGoogleOAuthToken(tok)
	t.Log("output:", output)
	if err != nil {
		t.Log("ERROR:", err)
	}
}

func TestGetGoogleOAuthToken(t *testing.T) {
	tok, err := testTableBasics.GetGoogleOAuthToken("test123")
	if err != nil {
		t.Log(err)
	}
	t.Log("Get Token:", tok)
}
