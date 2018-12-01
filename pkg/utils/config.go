package utils

import (
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type config struct {
	*viper.Viper
}

// Config set/get app vars
var Config = &config{viper.New()}

func init() {
	//jump to root directory
	_, fileName, _, _ := runtime.Caller(0)
	rootPath := path.Join(fileName, "../../../configs/")
	err := os.Chdir(rootPath)
	if err != nil {
		panic(err)
	}
	Config.AddConfigPath(".")
	Config.SetConfigName("config")
	err = Config.ReadInConfig()
	if err != nil {
		Log.Errorf("Fatal error config file: %s", err)
	}
}

func (c *config) GetString(key string) string {
	//search environment variable, if not found, then configuration file
	//environment variable is capitalized string, replace '.' in key as '_'
	KEY := strings.Replace(key, ".", "_", -1)
	value := os.Getenv(strings.ToUpper(KEY))
	if value == "" {
		//not found
		value = c.Viper.GetString(key)
	}
	return value
}
