package models

import (
	"fmt"
	"github.com/evoself/spiderhub/pkg/db/mongodb"
	"testing"
)

func TestPager_Search(t *testing.T) {
	mgo.Open()
	p := NewPage()
	fmt.Println(p.Search("盗梦空间"))
}
