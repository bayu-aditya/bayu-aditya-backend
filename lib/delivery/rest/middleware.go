package rest

import (
	"github.com/gin-gonic/gin"
)

// CORS for website
func CORS() gin.HandlerFunc {
	return func(c *gin.Context) {
		origin := c.Request.Header.Get("origin")
		switch origin {
		case
			"http://localhost:3000",          // localhost frontend
			"https://localhost:3000",         // localhost frontend
			"https://staging.bayuaditya.dev", // staging
			"https://bayuaditya.dev",         // production
			"https://www.bayuaditya.dev":     // production

			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
		}

		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Disposition, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, PATCH, DELETE")
		c.Writer.Header().Set("Access-Control-Max-Age", "600")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
