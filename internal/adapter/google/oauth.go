package google

import (
	"context"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleOAuth struct {
	Config *oauth2.Config
}

func NewGoogleOAuth(clientID string, clientSecret string, redirectURL string) *GoogleOAuth {
	return &GoogleOAuth{
		Config: &oauth2.Config{
			ClientID:     clientID,
			ClientSecret: clientSecret,
			Endpoint:     google.Endpoint,
			Scopes:       []string{drive.DriveScope},
			RedirectURL:  redirectURL,
		},
	}
}

func (oa *GoogleOAuth) OAuthLoginURL(lineID string) (oauthURL string) {
	oauthURL = oa.Config.AuthCodeURL(lineID, oauth2.AccessTypeOffline, oauth2.ApprovalForce) // oauth2.ApprovalForce
	return oauthURL
}

func (oa *GoogleOAuth) UserOAuthToken(authCode string) (*oauth2.Token, error) {
	token, err := oa.Config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Printf("Unable to retrieve token from web %v", err)
		return nil, err
	}
	return token, nil
}

func (oa *GoogleOAuth) NewGoogleDrive(ctx context.Context, tok *oauth2.Token) (*GoogleDrive, error) {
	client := oa.Config.Client(context.Background(), tok)
	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Unable to retrieve Drive client: %v", err)
		return nil, err
	}
	return &GoogleDrive{
		Service: srv,
	}, nil
}
