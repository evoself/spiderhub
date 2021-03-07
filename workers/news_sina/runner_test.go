package news_sina

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
