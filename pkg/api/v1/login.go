package v1

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func login(c gin.Context) {
	session := sessions.Default(c)

}
