package models

import (
	"fmt"
	"spiderhub/pkg/mgo"
	"testing"
)

func TestPager_Search(t *testing.T) {
	mgo.Open()
	p := NewPage()
	fmt.Println(p.Search("盗梦空间"))
}
