package utils

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"time"
)

var (
	// Redis - Global Redis Client Ref.
	RedisClient *redis.Client
	// KeyPrefix key prefix
	KeyPrefix string
)

func init() {
	// init redis client with options
	options := redis.Options{
		Addr:     Config.GetString("cache.redis.host") + ":" + Config.GetString("cache.redis.port"),
		Password: Config.GetString("cache.redis.password")}
	RedisClient = redis.NewClient(&options)
	KeyPrefix = Config.GetString("cache.redis.prefix")
}

func RedisSet(key string, value string, expiration time.Duration) error {
	err := RedisClient.Set(key, value, expiration).Err()
	if err != nil {
		// redis连接失败等
		Log.Errorf("RedisSet fail, error: [%v] ", err)
		return err
	} else {
		return nil
	}
}

// 测试key对应的值是不是expect
func RedisVerifyValue(key string, val string) error {
	expect, err := RedisClient.Get(key).Result()
	if err == redis.Nil {
		Log.Errorf("key [%s] does not exist in redis", key)
		return errors.New("key does not exist in redis")
	} else if err != nil {
		// redis连接失败等
		Log.Errorf("RedisVerifyValue fail, error: [%v] ", err)
		return err
	} else {
		if val != expect {
			// 找到了，但是不一致
			msg := fmt.Sprintf("expect %s, but got %s", expect, val)
			Log.Errorf("RedisVerifyValue fail, error: [%v] ", msg)
			return errors.New(msg)
		} else {
			// 找到了，并且一致
			return nil
		}
	}
}

func GetCacheSetMembers(key string) ([]string,error) {
	all, err := RedisClient.SMembers(key).Result()
	if err != nil {
		Log.Warnf("Get set objects failed: %v", err)
		return all,err
	}
	return all,nil
}

func SetCacheSetMember(key string,member interface{}) {
	if err := RedisClient.SAdd(key,member).Err();err != nil {
		Log.Warnf("Error in caching set of objects: %v", err)
	}
}

func DelCacheSetMember(key string,member interface{}) {
	if err := RedisClient.SRem(key,member).Err();err != nil {
		Log.Warnf("Error in caching set of objects: %v", err)
	}
}

func UniqueMerchantOnlineAutoKey() string {
	return KeyPrefix + ":merchant:online:auto:list"
}

func UniqueMerchantOnlineKey() string {
	return KeyPrefix + ":merchant:online:auto"
}
