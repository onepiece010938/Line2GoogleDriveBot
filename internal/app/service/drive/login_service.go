package drive

import (
	"context"
	"log"

	"github.com/onepiece010938/Line2GoogleDriveBot/internal/adapter/dynamodb"
	domainDrive "github.com/onepiece010938/Line2GoogleDriveBot/internal/domain/drive"
)

func (dr *GoogleDriveService) Login(ctx context.Context, lineID string, authCode string) error {
	tok, err := dr.driveServiceGoogleOA.UserOAuthToken(authCode)
	if err != nil {
		return err
	}
	// fmt.Println("Login service GET", lineID)
	// fmt.Printf("%#v", dr.driveServiceDynamodb)

	dToken := dynamodb.GoogleOAuthToken{
		PK:           lineID,
		AccessToken:  tok.AccessToken,
		TokenType:    tok.TokenType,
		RefreshToken: tok.RefreshToken,
		Expiry:       tok.Expiry,
		Info: map[string]interface{}{
			"upload_folder_id": ""},
	}

	err = dr.driveServiceDynamodb.AddGoogleOAuthToken(dToken)
	if err != nil {
		return err
	}

	// fmt.Printf("%#v", tok)
	// fmt.Printf("%#v", dToken)
	return nil
}

func (dr *GoogleDriveService) LoginURL(ctx context.Context, lineID string) string {
	oauthURL := dr.driveServiceGoogleOA.OAuthLoginURL(lineID)
	resURL, err := domainDrive.AppendOpenExternalBrowserParam(oauthURL)
	if err != nil {
		log.Println("Error:", err)
		return ""
	}
	return resURL
}

/*
Login

check token
if expored update token
save to db
*/

/*
Save token
*/
