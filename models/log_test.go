package models

import (
	"spiderhub/pkg/mgo"
	"testing"
)

func Test(t *testing.T) {
	mgo.Open()
	_, _ = NewLog("err!").Save()
}
