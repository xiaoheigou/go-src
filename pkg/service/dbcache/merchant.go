package dbcache

import (
	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"time"
	"yuudidi.com/pkg/utils"
)

// 首先从redis中读取记录，如果读不到，则到数据库中读取，并加载在到redis中
func getRecordById(cachePrefix string, id int64, dataPointer interface{}) error {
	var expire = time.Duration(10) * time.Minute  // 缓存10分钟过期
	var idStr = strconv.FormatInt(id, 10)
	var key = cachePrefix + idStr

	var redisData string
	var err error
	redisData, err = utils.RedisClient.Get(key).Result()
	if err == redis.Nil {
		utils.Log.Infof("id=[%d] does not exist in redis (prefix=[%s]), we would load it into redis later", id, cachePrefix)
	} else if err != nil {
		// redis连接失败等
		utils.Log.Errorf("redis fail, error: [%v] ", err)
	} else {
		if err = bson.Unmarshal([]byte(redisData), dataPointer); err != nil {
			utils.Log.Errorf("redis unmarshal uid=[%d] fail (prefix=[%s]), error: [%v] ", id, cachePrefix, err)
		} else {
			// redis中找到数据，且unmarshal成功，则返回数据
			return nil
		}
	}

	// Cache中没找到，读数据库
	if err := utils.DB.Where("id = ?", id).Find(dataPointer).Error; err != nil {
		utils.Log.Errorf("find id=[%d] fail (prefix=[%s]). [%v]", id, cachePrefix, err)
		return err
	}

	// 把数据转为bson后存入redis缓存
	bsonData, _ := bson.Marshal(dataPointer)
	if err := utils.RedisClient.Set(key, bsonData, expire).Err(); err != nil {
		utils.Log.Errorf("save id=[%d] to redis fail (prefix=[%s]). [%v]", id, cachePrefix, err)
	}

	return nil
}

func deleteKeyInRedis(cachePrefix string, id int64) error {
	var idStr = strconv.FormatInt(id, 10)
	var key = cachePrefix + idStr
	if err := utils.RedisClient.Del(key).Err(); err != nil {
		utils.Log.Errorf("delete key [%d] in redis fail. err = [%v]", key, err)
		return err
	}
	return nil
}

func GetMerchantById(id int64, dataPointer interface{})  error {
	var prefix = "db:table:merchants:"
	if err := getRecordById(prefix, id, dataPointer); err != nil {
		utils.Log.Warnln("find merchant id=%d fail. err = [%v]", err)
		return err
	}
	return nil
}

// 修改了merchants表后，请务必调用这个函数使对应的缓存失效
func InvalidateMerchant(id int64) error {
	var prefix = "db:table:merchants:"
	return deleteKeyInRedis(prefix, id)
}

func GetPreferenceById(id int64, dataPointer interface{})  error {
	var prefix = "db:table:preferences:"
	if err := getRecordById(prefix, id, dataPointer); err != nil {
		utils.Log.Warnln("find preference id=%d fail. err = [%v]", err)
		return err
	}
	return nil
}

// 修改了perferences表后，请务必调用这个函数使对应的缓存失效
func InvalidatePreference(id int64) error {
	var prefix = "db:table:preferences:"
	return deleteKeyInRedis(prefix, id)
}

//func main() {
//	var perf models.Preferences
//	var err error
//	if err = GetPreferenceById(1, &perf); err != nil {
//		utils.Log.Warnf("err [%v]", err)
//	}
//	fmt.Println(perf.AutoConfirm)
//}