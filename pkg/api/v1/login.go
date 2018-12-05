package v1

import (
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {

}

func AppLogin(c *gin.Context) {
	//TODO

	c.JSON(200, gin.H{
		"token": 123,
		"uid":   123,
	})
}


