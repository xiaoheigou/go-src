package middleware

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/utils"
)

// Authenticated - check if user authenticated
func Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("123")
		utils.Log.Infof("username:%v",username)
		if username == nil {
			c.AbortWithError(401, errors.New("Access Forbidden"))
		}
	}
}