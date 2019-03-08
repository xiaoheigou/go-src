package utils

import (
	"errors"
	"fmt"
	"github.com/go-redis/redis"
	"strconv"
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

func GetCacheSetMembers(key string) ([]string, error) {
	all, err := RedisClient.SMembers(key).Result()
	if err != nil {
		Log.Warnf("Get set objects failed: %v", err)
		return all, err
	}
	return all, nil
}

// 得到多个集合中元素的交集，结果保存在result中
func GetCacheSetInterMembers(result *[]string, keys ...string) error {
	all, err := RedisClient.SInter(keys...).Result()
	*result = all
	if err != nil {
		Log.Warnf("Get set objects failed: %v", err)
		return err
	}
	return nil
}

func SetCacheSetMember(key string, expireTime int, member ...interface{}) error {
	if err := RedisClient.SAdd(key, member...).Err(); err != nil {
		Log.Warnf("Error in caching set of objects: %v", err)
		return err
	}
	if expireTime > 0 {
		if err := RedisClient.Expire(key, time.Duration(expireTime)*time.Second).Err(); err != nil {
			Log.Warnf("Error in set expire time,key:%s", key)
		}
	}
	return nil
}

func DelCacheSetMember(key string, member ...interface{}) error {
	if err := RedisClient.SRem(key, member...).Err(); err != nil {
		Log.Warnf("Error in caching set of objects: %v", err)
		return err
	}
	return nil
}

func UpdateMerchantLastOrderTime(merchantId int64, direction int, acceptAt time.Time) error {
	score := float64(acceptAt.Unix())

	var key string
	if direction == 0 {
		key = RedisKeyMerchantLastAcceptTimeForD0Order()
	} else if direction == 1 {
		key = RedisKeyMerchantLastAcceptTimeForD1Order()
	} else {
		return errors.New("invalid param direction")
	}

	if err := RedisClient.ZAdd(key, redis.Z{Score: score, Member: strconv.FormatInt(merchantId, 10)}).Err(); err != nil {
		Log.Warnf("redis zadd error: %v", err)
		return err
	}
	return nil
}

func GetMerchantsSortedByLastOrderTime(direction int) ([]string, error) {
	var key string
	if direction == 0 {
		key = RedisKeyMerchantLastAcceptTimeForD0Order()
	} else if direction == 1 {
		key = RedisKeyMerchantLastAcceptTimeForD1Order()
	} else {
		return []string{}, errors.New("invalid param direction")
	}

	var sortedResult []string
	var err error
	if sortedResult, err = RedisClient.ZRangeByScore(key, redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result(); err != nil {
		Log.Warnf("redis zrangebyscore error: %v", err)
		return []string{}, err
	}
	return sortedResult, nil
}

// 保存自动订单的派单时间
func UpdateMerchantLastAutoOrderSendTime(merchantId int64, direction int, sendAt time.Time) error {
	score := float64(sendAt.Unix())

	var key string
	if direction == 0 {
		key = RedisKeyMerchantLastSendTimeForD0AutoOrder()
	} else {
		// 目前只有用户充值订单（direction == 0）会是自动订单。
		return errors.New("invalid param direction, only 0 is expected for auto order")
	}

	if err := RedisClient.ZAdd(key, redis.Z{Score: score, Member: strconv.FormatInt(merchantId, 10)}).Err(); err != nil {
		Log.Warnf("redis zadd error: %v", err)
		return err
	}
	return nil
}

// 返回币商列表，按自动订单的派单时间排序
func GetMerchantsSortedByLastAutoOrderSendTime(direction int) ([]string, error) {
	var key string
	if direction == 0 {
		key = RedisKeyMerchantLastSendTimeForD0AutoOrder()
	} else {
		return []string{}, errors.New("invalid param direction, only 0 is expected for auto order")
	}

	var sortedResult []string
	var err error
	if sortedResult, err = RedisClient.ZRangeByScore(key, redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
	}).Result(); err != nil {
		Log.Warnf("redis zrangebyscore error: %v", err)
		return []string{}, err
	}
	return sortedResult, nil
}

func IncreaseAppLoginFailTimes(nationCode int, phone string) error {
	key := KeyPrefix + ":app:loginfail:" + strconv.Itoa(nationCode) + ":" + phone

	_, err := RedisClient.Get(key).Result()
	if err == redis.Nil {
		// 没找到记录
		var lockHours int64
		if lockHours, err = strconv.ParseInt(Config.GetString("app.loginfail.lockhours"), 10, 0); err != nil {
			Log.Errorf("Wrong configuration: app.loginfail.lockhours, should be int. Set to default 24.")
			lockHours = 24
		}
		expiration := time.Duration(lockHours) * time.Hour

		Log.Debugf("key [%s] does not exist in redis", key)
		if err1 := RedisClient.Set(key, 1, expiration).Err(); err1 != nil {
			Log.Errorf("RedisClient.Set fail, error: [%v] ", err1)
			return err1
		}
		return nil
	} else if err != nil {
		// redis连接失败等
		Log.Errorf("IncreaseAppLoginFailTimes fail, error: [%v] ", err)
		return err
	} else {
		// 找到记录，增加次数
		if err1 := RedisClient.Incr(key).Err(); err1 != nil {
			Log.Errorf("RedisClient.Incr fail, error: [%v] ", err1)
			return err1
		}
		return nil
	}
}

func ReachMaxAppLoginFailTimes(nationCode int, phone string) bool {
	key := KeyPrefix + ":app:loginfail:" + strconv.Itoa(nationCode) + ":" + phone

	got, err := RedisClient.Get(key).Result()
	if err == redis.Nil {
		// 没找到记录
		return false
	} else if err != nil {
		// redis连接失败等
		Log.Errorf("GetAppLoginFailTimes fail, error: [%v] ", err)
		return false
	} else {
		// 找到记录
		var gotInt int
		var err1 error
		if gotInt, err1 = strconv.Atoi(got); err1 != nil {
			return false
		}

		var maxTimes int64
		if maxTimes, err = strconv.ParseInt(Config.GetString("app.loginfail.maxtimes"), 10, 0); err != nil {
			Log.Errorf("Wrong configuration: app.loginfail.maxtimes, should be int. Set to default 3.")
			maxTimes = 3
		}

		if gotInt >= int(maxTimes) {
			// 记录的失败次数太多
			Log.Errorf("%d %s login fail %d times", nationCode, phone, gotInt)
			return true
		}
		return false
	}
}

func ClearAppLoginFailTimes(nationCode int, phone string) error {
	key := KeyPrefix + ":app:loginfail:" + strconv.Itoa(nationCode) + ":" + phone

	if err := RedisClient.Del(key).Err(); err != nil {
		Log.Errorf("delete key [%s] in redis fail. err = [%v]", key, err)
		return err
	}

	return nil
}

func RedisIncreaseRefulfillTimesToOfficialMerchants(orderNumber string) error {
	key := KeyPrefix + ":order:officialmerchants:trytimes:" + orderNumber

	_, err := RedisClient.Get(key).Result()
	if err == redis.Nil {
		// 没找到记录
		Log.Debugf("key [%s] does not exist in redis", key)
		if err1 := RedisClient.Set(key, 1, 0).Err(); err1 != nil {
			Log.Errorf("RedisClient.Set fail, error: [%v] ", err1)
			return err1
		}
		return nil
	} else if err != nil {
		// redis连接失败等
		Log.Errorf("RedisIncreaseRefulfillTimesToOfficialMerchants fail, error: [%v] ", err)
		return err
	} else {
		// 找到记录，增加次数
		if err1 := RedisClient.Incr(key).Err(); err1 != nil {
			Log.Errorf("RedisClient.Incr fail, error: [%v] ", err1)
			return err1
		}
		return nil
	}
}

// 如果失败返回0
func RedisGetRefulfillTimesToOfficialMerchants(orderNumber string) int64 {
	key := KeyPrefix + ":order:officialmerchants:trytimes:" + orderNumber

	var ret int64 = 0
	got, err := RedisClient.Get(key).Result()
	if err == redis.Nil {
		// 没找到记录
		return ret
	} else if err != nil {
		// redis连接失败等
		Log.Errorf("RedisGetRefulfillTimesToOfficialMerchants fail, error: [%v] ", err)
		return ret
	} else {
		// 找到记录
		var gotInt int64
		var err1 error
		if gotInt, err1 = strconv.ParseInt(got, 10, 64); err1 != nil {
			return 0
		}

		return gotInt
	}
}

func RedisDelRefulfillTimesToOfficialMerchants(orderNumber string) error {
	key := KeyPrefix + ":order:officialmerchants:trytimes:" + orderNumber

	if err := RedisClient.Del(key).Err(); err != nil {
		Log.Errorf("delete key [%s] in redis fail. err = [%v]", key, err)
		return err
	}

	return nil
}

func RedisIncreaseRefulfillTimesToNormalMerchants(orderNumber string) error {
	key := KeyPrefix + ":order:normalmerchants:trytimes:" + orderNumber

	_, err := RedisClient.Get(key).Result()
	if err == redis.Nil {
		// 没找到记录
		Log.Debugf("key [%s] does not exist in redis", key)
		if err1 := RedisClient.Set(key, 1, 0).Err(); err1 != nil {
			Log.Errorf("RedisClient.Set fail, error: [%v] ", err1)
			return err1
		}
		return nil
	} else if err != nil {
		// redis连接失败等
		Log.Errorf("RedisIncreaseRefulfillTimesToNormalMerchants fail, error: [%v] ", err)
		return err
	} else {
		// 找到记录，增加次数
		if err1 := RedisClient.Incr(key).Err(); err1 != nil {
			Log.Errorf("RedisClient.Incr fail, error: [%v] ", err1)
			return err1
		}
		return nil
	}
}

// 如果失败返回0
func RedisGetRefulfillTimesToNormalMerchants(orderNumber string) int64 {
	key := KeyPrefix + ":order:normalmerchants:trytimes:" + orderNumber

	var ret int64 = 0
	got, err := RedisClient.Get(key).Result()
	if err == redis.Nil {
		// 没找到记录
		return ret
	} else if err != nil {
		// redis连接失败等
		Log.Errorf("RedisGetRefulfillTimesToNormalMerchants fail, error: [%v] ", err)
		return ret
	} else {
		// 找到记录
		var gotInt int64
		var err1 error
		if gotInt, err1 = strconv.ParseInt(got, 10, 64); err1 != nil {
			return 0
		}

		return gotInt
	}
}

func RedisDelRefulfillTimesToNormalMerchants(orderNumber string) error {
	key := KeyPrefix + ":order:normalmerchants:trytimes:" + orderNumber

	if err := RedisClient.Del(key).Err(); err != nil {
		Log.Errorf("delete key [%s] in redis fail. err = [%v]", key, err)
		return err
	}

	return nil
}

func RedisKeyMerchantRole1() string {
	return KeyPrefix + ":merchant:role1"
}

func RedisKeyMerchantOnline() string {
	return KeyPrefix + ":merchant:online"
}

func RedisKeyMerchantInWork() string {
	return KeyPrefix + ":merchant:in_work"
}

func RedisKeyMerchantWechatHookStatus() string {
	return KeyPrefix + ":merchant:wechat_hook_status"
}

func RedisKeyMerchantAlipayHookStatus() string {
	return KeyPrefix + ":merchant:alipay_hook_status"
}

func RedisKeyMerchantWechatAutoOrder() string {
	return KeyPrefix + ":merchant:wechat_auto_order"
}

func RedisKeyMerchantAlipayAutoOrder() string {
	return KeyPrefix + ":merchant:alipay_auto_order"
}

func RedisKeyMerchantSelected(orderNumber string) string {
	return KeyPrefix + ":merchant:selected:" + orderNumber
}

func RedisKeyMerchantLastAcceptTimeForD0Order() string {
	return KeyPrefix + ":merchant:direction_0_last_order_time" // 记录币商最近一次direction 0的订单的接单时间
}

func RedisKeyMerchantLastAcceptTimeForD1Order() string {
	return KeyPrefix + ":merchant:direction_1_last_order_time" // 记录向商最近一次direction 1的订单的接单时间
}

func RedisKeyMerchantLastSendTimeForD0AutoOrder() string {
	return KeyPrefix + ":merchant:direction_0_auto_order_last_send_time" // 记录币商最近一次direction 0的自动订单的派单时间
}

func UniqueDistributorTokenKey(token string) string {
	return KeyPrefix + ":distributor:" + token
}

func UniqueTimeWheelKey(sign string) string {
	return KeyPrefix + ":timewheel:" + sign
}
