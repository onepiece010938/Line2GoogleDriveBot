package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		c.Header("Access-Control-Allow-Origin", "*")                                                     // "*" can change to certain domain
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization") //add the header that u accept
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type,Content-Disposition,Content-Description,Content-Transfer-Encoding")
		c.Header("Access-Control-Allow-Credentials", "true")
		//c.Header("Transfer-Encoding", "chunked")
		//c.Header("Content-Transfer-Encoding", "chunked")
		// accept all OPTIONS
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}
