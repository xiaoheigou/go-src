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
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		Config.GetString("database.user"),
		Config.GetString("database.pass"),
		Config.GetString("database.host"),
		Config.GetString("database.port"),
		Config.GetString("database.name")))
	if err != nil {
		Log.Error("database connect error: ", err)
	}

	DB.LogMode(Config.GetBool("database.debug"))
}
