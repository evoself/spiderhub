package news_tencent

import (
	"spiderhub/config"
	"spiderhub/pkg/mgo"
	"testing"
)

func Test(t *testing.T) {
	config.Load()
	mgo.Open()
	Run()
}