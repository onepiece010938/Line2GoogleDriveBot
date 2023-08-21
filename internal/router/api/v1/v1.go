package v1

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/onepiece010938/Line2GoogleDriveBot/internal/app"
	"github.com/tidwall/gjson"

	"github.com/gin-gonic/gin"
	// "golang.org/x/oauth2"
	// "golang.org/x/oauth2/google"
	// drive "google.golang.org/api/drive/v3"
)

const redirectURL = "https://bf62-220-134-22-223.ngrok-free.app/api/v1/ouath_login"
const scope = "https://www.googleapis.com/auth/drive"

// https://www.googleapis.com/auth/drive
// "https://www.googleapis.com/auth/userinfo.profile"

// Config Config
type Config struct {
	GoogleSecretKey string `mapstructure:"GOOGLE_SECRET_KEY"`
	GoogleClientID  string `mapstructure:"GOOLE_CLIENT_ID"`
}

var Val Config

/*
func oauthURL2() string {
	// ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     "545014009419-em4rm9dqmd0lug3nlsv68vt3eafugbdo.apps.googleusercontent.com", // from https://console.developers.google.com/project/<your-project-id>/apiui/credential
		ClientSecret: "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",                                      // from https://console.developers.google.com/project/<your-project-id>/apiui/credential
		Endpoint:     google.Endpoint,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.profile", drive.DriveScope},
	}
	// Redirect user to consent page to ask for permission
	// for the scopes specified above.
	randState := fmt.Sprintf("st%d", time.Now().UnixNano())
	url := conf.AuthCodeURL(randState)
	fmt.Printf("Visit the URL for the auth dialog: %v", url)

	// var code string
	// if _, err := fmt.Scan(&code); err != nil {
	// 	log.Fatal(err)
	// }
	// tok, err := conf.Exchange(ctx, code)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// client := conf.Client(ctx, tok)
	// client.Get("...")

	return url
}
*/

func RegisterRouter(router *gin.RouterGroup, app *app.Application) {
	OauthInit()
	Val = Config{
		GoogleSecretKey: "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",
		GoogleClientID:  "545014009419-em4rm9dqmd0lug3nlsv68vt3eafugbdo.apps.googleusercontent.com",
	}

	v1 := router.Group("/v1")
	{
		v1.POST("/callback", Callback(app))
		v1.GET("/sample", SAMPLE)
		v1.GET("/analyze", StartAnalyze(app))
		v1.GET("/ouath_url", GoogleAccsessNew)
		v1.GET("/ouath_login", GoogleLoginNew)
		v1.GET("/drive", GetDrive)

		v1.GET("/ouath_url2", GoogleAccsess)
		v1.GET("/ouath_login2", GoogleLogin)
		v1.GET("/xx", Main2)
	}
}

func GoogleAccsess(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, gin.H{
		"url": oauthURL(),
	})
}

func oauthURL() string {
	u := "https://accounts.google.com/o/oauth2/v2/auth?client_id=%s&response_type=code&scope=%s&redirect_uri=%s"

	return fmt.Sprintf(u, Val.GoogleClientID, scope, redirectURL)
}

// GoogleLogin GoogleLogin
func GoogleLogin(c *gin.Context) {
	code := c.Query("code")

	token, err := accessToken(code)
	if err != nil {
		log.Println("err", "accessToken error: ", err)
		c.Redirect(http.StatusFound, "/")
		return
	}

	id, name, err := getGoogleUserInfo(token)
	if err != nil {
		log.Println("err", "getGoogleUserInfo error: ", err)
		c.Redirect(http.StatusFound, "/")
		return
	}

	log.Printf("id: %v, name: %v", id, name)
}

func accessToken(code string) (token string, err error) {
	u := "https://www.googleapis.com/oauth2/v4/token"

	data := url.Values{"code": {code}, "client_id": {Val.GoogleClientID}, "client_secret": {Val.GoogleSecretKey}, "grant_type": {"authorization_code"}, "redirect_uri": {redirectURL}}
	body := strings.NewReader(data.Encode())

	resp, err := http.Post(u, "application/x-www-form-urlencoded", body)
	if err != nil {
		return token, err
	}

	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return token, err
	}

	token = gjson.GetBytes(b, "access_token").String()

	return token, nil
}

func getGoogleUserInfo(token string) (id, name string, err error) {
	u := fmt.Sprintf("https://www.googleapis.com/oauth2/v1/userinfo?alt=json&access_token=%s", token)
	resp, err := http.Get(u)
	if err != nil {
		return id, name, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return id, name, err
	}

	name = gjson.GetBytes(body, "name").String()
	id = gjson.GetBytes(body, "id").String()

	return id, name, nil
}
