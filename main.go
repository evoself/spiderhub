package main

import (
	"github.com/spf13/viper"
	"log"
	"spiderhub/pkg/mgo"
	"spiderhub/workers/news_163"
	"spiderhub/workers/news_tencent"
	"sync"
)

// init 初始化配置
func init() {
	load()
	mgo.Open()
}

func load() {
	viper.SetConfigFile("config.toml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		news_tencent.Run()
		wg.Done()
	}()
	go func() {
		news_163.Run()
		wg.Done()
	}()
	wg.Wait()
}
