package v1

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

// Retrieve a token, saves the token, then returns the generated client.
func getClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "token.json"
	tok, err := tokenFromFile(tokFile)
	if err != nil {
		tok = getTokenFromWeb(config)
		saveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}

// Request a token from the web, then returns the retrieved token.
func getTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	var authCode string
	if _, err := fmt.Scan(&authCode); err != nil {
		log.Fatalf("Unable to read authorization code %v", err)
	}

	tok, err := config.Exchange(context.TODO(), authCode)
	if err != nil {
		log.Fatalf("Unable to retrieve token from web %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func tokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func saveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()
	err = json.NewEncoder(f).Encode(token)
	if err != nil {
		log.Printf("Unable to Encode oauth token: %v", err)
	}
}

func Main2(c *gin.Context) {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	// If modifying these scopes, delete your previously saved token.json.
	config, err := google.ConfigFromJSON(b, drive.DriveScope) //DriveMetadataReadonlyScope
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}
	client := getClient(config)

	srv, err := drive.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Fatalf("Unable to retrieve Drive client: %v", err)
	}

	r, err := srv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
}

var CACHE_TOKEN *oauth2.Token

func GetDrive(c *gin.Context) {
	client := CONFIGG.Client(context.Background(), CACHE_TOKEN)
	srv, err := drive.NewService(c, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Unable to retrieve Drive client: %v", err)
	}

	r, err := srv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Printf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
}

func GoogleLoginNew(c *gin.Context) {
	authCode := c.Query("code")

	// tokFile := "token.json"
	// tok, err := tokenFromFile(tokFile)
	// if err != nil {
	// 	// tok, _ := CONFIGG.Exchange(context.TODO(), authCode)
	// 	saveToken(tokFile, tok)
	// }
	tok, err := CONFIGG.Exchange(context.TODO(), authCode)
	CACHE_TOKEN = tok
	saveToken("token.json", CACHE_TOKEN)
	log.Println("@@TOKEN: ", tok)
	if err != nil {
		log.Printf("Unable to retrieve token from web %v", err)
	}
	client := CONFIGG.Client(context.Background(), tok)
	// driveService, err := drive.NewService(ctx, option.WithTokenSource(CONFIGG.TokenSource(c, token)))
	srv, err := drive.NewService(c, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("Unable to retrieve Drive client: %v", err)
	}

	r, err := srv.Files.List().PageSize(10).
		Fields("nextPageToken, files(id, name)").Do()
	if err != nil {
		log.Printf("Unable to retrieve files: %v", err)
	}
	fmt.Println("Files:")
	if len(r.Files) == 0 {
		fmt.Println("No files found.")
	} else {
		for _, i := range r.Files {
			fmt.Printf("%s (%s)\n", i.Name, i.Id)
		}
	}
	c.Writer.Write([]byte("<html><title>BaoSave Login</title> <body> Authorized successfully, please close this window</body></html>"))

	// filename := `C:\Users\raymond\Desktop\MyGitRepo\Line2GoogleDriveBot\README.md`
	// goFile, err := os.Open(filename)
	// if err != nil {
	// 	log.Fatalf("error opening %q: %v", filename, err)
	// }
	// driveFile, err := srv.Files.Create(&drive.File{Name: filename}).Media(goFile).Do()
	// log.Printf("Got drive.File, err: %#v, %v", driveFile, err)
}

func GoogleAccsessNew(ctx *gin.Context) {

	authURL := CONFIGG.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	fmt.Printf("Go to the following link in your browser then type the "+
		"authorization code: \n%v\n", authURL)

	url := fmt.Sprintf("%v", authURL)

	ctx.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

var CONFIGG *oauth2.Config

func OauthInit() {
	// b, err := os.ReadFile("credentials.json")
	// if err != nil {
	// 	log.Fatalf("Unable to read client secret file: %v", err)
	// }

	// // If modifying these scopes, delete your previously saved token.json.
	// CONFIGG, err = google.ConfigFromJSON(b, drive.DriveScope) //DriveMetadataReadonlyScope
	// if err != nil {
	// 	log.Fatalf("Unable to parse client secret file to config: %v", err)
	// }

	CONFIGG = &oauth2.Config{
		ClientID:     os.Getenv("ClientID"),
		ClientSecret: os.Getenv("ClientSecret"), // from https://console.developers.google.com/project/<your-project-id>/apiui/credential
		Endpoint:     google.Endpoint,
		Scopes:       []string{drive.DriveScope},
		RedirectURL:  os.Getenv("RedirectURL"),
	}

	tokFile := "token.json"
	tok, _ := tokenFromFile(tokFile)
	CACHE_TOKEN = tok

}
