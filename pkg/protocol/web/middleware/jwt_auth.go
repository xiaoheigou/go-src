package middleware

import (
	"errors"
	jwt_lib "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gin-gonic/gin"
	"yuudidi.com/pkg/utils"
)

func Auth(secret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := request.ParseFromRequest(c.Request, request.OAuth2Extractor, func(token *jwt_lib.Token) (interface{}, error) {
			b := ([]byte(secret))
			return b, nil
		})
		if err != nil {
			utils.Log.Errorf("Authorization fail [%v]", err)
			c.AbortWithError(401, err)
			return
		}

		if claims, ok := token.Claims.(jwt_lib.MapClaims); ok && token.Valid {
			tokenUid := claims["uid"]
			resourceUid := c.Param("uid")
			if tokenUid != resourceUid {
				utils.Log.Errorf("jwt can only access resource belong to uid [%v], but you want to access resource belong to uid [%s]", tokenUid, resourceUid)
				c.AbortWithError(401, errors.New("Authorization fail"))
				return
			}
		} else {
			utils.Log.Errorln("Parse jwt error")
			c.AbortWithError(401, errors.New("Parse jwt error"))
			return
		}
	}
}
