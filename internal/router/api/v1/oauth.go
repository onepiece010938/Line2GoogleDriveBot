package v1

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/onepiece010938/Line2GoogleDriveBot/internal/app"
)

func OAuthLogin(app *app.Application) gin.HandlerFunc {
	return func(c *gin.Context) {
		authCode := c.Query("code")
		lineID := c.Query("state") //get lineid
		err := app.DriveService.Login(c.Request.Context(), lineID, authCode)
		if err != nil {
			log.Printf("Unable to retrieve DriveService.Login %v", err)
			c.String(http.StatusInternalServerError, "Unable to retrieve DriveService.Login")
			return
		}
		_, err = c.Writer.Write([]byte("<html><title>BaoSave Login</title> <body> Authorized successfully, please close this window</body></html>"))
		if err != nil {
			log.Printf("Unable to write HTML: %v", err)
		}
	}
}
