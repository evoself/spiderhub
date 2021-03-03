package models

import (
	"github.com/evoself/spiderhub/pkg/db/mongodb"
	"testing"
)

func Test(t *testing.T) {
	mgo.Open()
	_, _ = NewLog("err!").Save()
}
