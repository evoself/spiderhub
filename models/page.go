package models

import (
	"context"
	"github.com/qiniu/qmgo"
	"go.mongodb.org/mongo-driver/bson"
	"spiderhub/pkg/mgo"
)

const pageCollection = "page"

type page struct {
	Id          int64  `bson:"id"`
	Title       string `bson:"title"`
	Url         string `bson:"url"`
	Hash        string `bson:"hash"`
	Image       string `bson:"image"`
	Summary     string `bson:"summary"`
	Source      string `bson:"source"`
	Category    string `bson:"category"`
	SubCategory string `bson:"subcategory"`
	Date        int    `bson:"date"`
	PublishTime string `bson:"publish_time"`
	ExtractTime string `bson:"extract_time"`
}

func NewPage() *page {
	return &page{}
}

func (p *page) FindOne() {
	_ = mgo.Init(pageCollection).Find(context.Background(), bson.M{"url": p.Url}).Limit(1).One(&p)
}

// 插入单个列表数据
func (p *page) Save() (*qmgo.InsertOneResult, error) {
	newPage := &page{}
	_ = mgo.Init(pageCollection).Find(context.Background(), bson.M{"url": p.Url}).Limit(1).One(&newPage)
	if (len(newPage.Url) == 0) && (len(p.Url) > 0) {
		return mgo.Init(pageCollection).InsertOne(context.Background(), p)
	}
	return nil, nil
}

// 1. 获取会员上一次的主题和页码
// 2. 根据主题和页码查找集合
func (p *page) Search(q string) []page {
	var (
		pages []page
		limit int64 = 8
	)
	mgo.Init(pageCollection).
		Find(context.Background(), bson.M{"title": bson.M{"$regex": q}}).
		Sort("-create_time").Limit(limit).All(&pages)
	return pages
}
