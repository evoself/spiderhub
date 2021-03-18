package game_17173

import (
	"spiderhub/configs"
	"spiderhub/pkg/mgo"
	"testing"
)

func Test(t *testing.T) {
	configs.Load()
	mgo.Open()
	Run()
}