package middleware

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// Authenticated - check if user authenticated
func Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		if userName, err := c.Cookie("username"); err != nil || userName == "" {
			c.AbortWithError(401, errors.New("Access Unauthorized"))
		}
	}
}