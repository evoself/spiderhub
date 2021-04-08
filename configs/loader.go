package configs

import (
	"github.com/spf13/viper"
	"log"
)

func Load() {
	viper.SetConfigFile("*") // *替换为配置文件的本地绝对路径
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}
