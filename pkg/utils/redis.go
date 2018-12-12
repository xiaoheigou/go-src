package utils

import (
	"fmt"
	"github.com/go-redis/redis"
	"time"
	"errors"
)

var RedisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})


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
func RedisVerifyValue(key string, expect string) error  {
	val, err := RedisClient.Get(key).Result()
	if err == redis.Nil {
		return errors.New("key does not exist")
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