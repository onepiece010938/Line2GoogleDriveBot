package middleware

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

/*For linebot*/
func Analyze() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		decoder := json.NewDecoder(ctx.Request.Body)
		fmt.Println(decoder)
		// err := decoder.Decode(&frontReq)
	}
}
