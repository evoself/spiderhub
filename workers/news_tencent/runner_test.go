package news_tencent

import (
	"github.com/evoself/spiderhub/pkg/db/mgo"
	"testing"
)

func Test(t *testing.T) {
	mgo.Open()
	Run()
}
