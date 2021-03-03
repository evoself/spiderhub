package news_tencent

import (
	"spiderhub/pkg/mgo"
	"testing"
)

func Test(t *testing.T) {
	mgo.Open()
	Run()
}
