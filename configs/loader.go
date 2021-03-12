package configs

import (
	"github.com/spf13/viper"
	"log"
)

// 本地绝对路径
func Load() {
	viper.SetConfigFile("/Users/evoself/code/go/spiderhub/configs/params.toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}