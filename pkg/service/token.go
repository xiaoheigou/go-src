package service

import (
	"github.com/go-redis/redis"
	"github.com/typa01/go-utils"
	"time"
	"yuudidi.com/pkg/utils"
)

func GenerateToken(distributorId string) string {
	token := tsgutils.GUID()

	timeout := utils.Config.GetInt64("token.timeout")
	if err := utils.RedisSet(utils.UniqueDistributorTokenKey(token), distributorId, time.Duration(timeout)*time.Second); err != nil {
		utils.Log.Errorf("set distributor token key is failed! distributorId = %s", distributorId)
		return ""
	}

	return token
}

func VerifyToken(token string) bool {

	if _, err := utils.RedisClient.Get(utils.UniqueDistributorTokenKey(token)).Result(); err == redis.Nil {
		utils.Log.Debugf("token key is not exist,token %s", token)
		return false
	}
	return true
}
