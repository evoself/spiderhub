package main

import (
	"spiderhub/config"
	"spiderhub/pkg/mgo"
	"spiderhub/workers/news_163"
	"spiderhub/workers/news_sina"
	"spiderhub/workers/news_sohu"
	"spiderhub/workers/news_tencent"
	"sync"
	"time"
)

func init() {
	config.Load()
	mgo.Open()
}

func delay(n time.Duration) {
	time.Sleep(time.Minute * n)
}

func schedule(f func(), n time.Duration) {
	ticker := time.NewTicker(time.Hour * n)
	for range ticker.C {
		f()
	}
}

func main() {
	wg := sync.WaitGroup{}
	wg.Add(4)
	go func() {
		delay(0)
		f := func() {
			news_tencent.Run()
			wg.Done()
		}
		schedule(f, 10)
	}()
	go func() {
		delay(10)
		f := func() {
			news_163.Run()
			wg.Done()
		}
		schedule(f, 10)
	}()
	go func() {
		delay(20)
		f := func() {
			news_sina.Run()
			wg.Done()
		}
		schedule(f, 10)
	}()
	go func() {
		delay(30)
		f := func() {
			news_sohu.Run()
			wg.Done()
		}
		schedule(f, 10)
	}()
	wg.Wait()
}
