package dynamodb

import "testing"

func TestCreateGoogleOAuthTable(t *testing.T) {
	tableDesc, err := testTableBasics.CreateGoogleOAuthTable()
	t.Log("tableDesc:", tableDesc)
	t.Log("ERROR:", err)
	// if err ==
	// var notFoundEx *types.ResourceNotFoundException
}

func TestAddGoogleOAuthToken(t *testing.T) {
	tok := GoogleOAuthToken{
		PK:           "test1234",
		AccessToken:  "test123",
		TokenType:    "Bearer",
		RefreshToken: "test123",
		Expiry:       "2023-08-22T11:31:54.2936004+08:00",
	}
	err := testTableBasics.AddGoogleOAuthToken(tok)
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
		PK:          "test1234",
		AccessToken: "test123456",
		// TokenType:    "Bearer",
		RefreshToken: "test123456",
		// Expiry:       "2023-08-22T11:31:54.2936004+08:00",
	}
	output, err := testTableBasics.TxUpdateGoogleOAuthToken(tok)
	t.Log("output:", output)
	if err != nil {
		t.Log("ERROR:", err)
	}
}

func TestGetGoogleOAuthToken(t *testing.T) {
	tok, err := testTableBasics.GetGoogleOAuthToken("test1234")
	if err != nil {
		t.Log(err)
	}
	t.Log("Get Token:", tok)
}
