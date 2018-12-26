package utils

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql" //blank import
)

var (
	// DB return gorm.DB pointer
	DB  *gorm.DB
	err error
)

func init() {
	DB, err = gorm.Open(Config.GetString("database.driver"), fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=UTC",
		Config.GetString("database.user"),
		Config.GetString("database.pass"),
		Config.GetString("database.host"),
		Config.GetString("database.port"),
		Config.GetString("database.name")))
	if err != nil {
		Log.Error("database connect error: ", err)
	}
	//设置数据库连接池
	maxIdle := Config.GetInt("database.maxidle")
	maxOpen := Config.GetInt("database.maxopen")
	DB.DB().SetMaxIdleConns(maxIdle)
	DB.DB().SetMaxOpenConns(maxOpen)

	DB.LogMode(Config.GetBool("database.debug"))
}
