package middleware

import (
	"errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"time"
	"yuudidi.com/pkg/utils"
)

var timeout = utils.Config.GetInt64("web.server.timeout")

// Authenticated - check if user authenticated
func Authenticated() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		username := session.Get("username")
		timestamp := session.Get("timestamp")
		utils.Log.Debugf("username:%v,timestamp:%d", username, timestamp)
		if username == nil {
			c.AbortWithError(401, errors.New("Access Forbidden"))
		}
		now := time.Now().Unix()
		diff := timestamp.(int64) - now
		if diff > 0 {
			if diff <= (timeout / 10) {
				utils.Log.Debugf("session will expire,username:%s", username)
				session.Set("timestamp", now+timeout)
				session.Save()
			}
		} else {
			utils.Log.Debugf("session already expire,username:%s", username)
			c.AbortWithError(401, errors.New("Access Forbidden"))
		}
	}
}
