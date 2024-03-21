package conf

import (
	"fmt"
	"github.com/spf13/viper"
)

func Init(filePath string) {
	if len(filePath) < 1 {
		filePath = "conf/app.yml"
	}
	viper.SetConfigFile(filePath)
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("parse app.yml config file fail: %s", err))
	}
}

func GetString(key string) string {
	return viper.GetString(key)
}

func GetBool(key string) bool {
	return viper.GetBool(key)
}

func GetInt64(key string) int64 {
	return viper.GetInt64(key)
}
