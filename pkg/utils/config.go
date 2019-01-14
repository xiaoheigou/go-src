package utils

import (
	"github.com/fsnotify/fsnotify"
	"os"
	"path"
	"regexp"
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

	// Watching and re-reading config files
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		Log.Infof("Config file changed: %s", e.Name)
	})
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
	if strings.Contains(value, "${") {
		value = replaceValue(value, c)
	}
	return value
}

func replaceValue(value string, c *config) string {
	regKey := regexp.MustCompile(`\$\{(.*?)\}`)
	values := regKey.FindAllStringSubmatch(value, -1)
	for _, v := range values {
		k := os.Getenv(strings.ToUpper(v[1]))
		if k == "" {
			k = c.Viper.GetString(v[1])
		}
		value = strings.Replace(value, v[0], k, 1)
	}
	return value
}
