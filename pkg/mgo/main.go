package mgo

import (
	"context"
	"fmt"
	"github.com/qiniu/qmgo"
	"github.com/spf13/viper"
)

var (
	ctx    context.Context
	client *qmgo.Client
)

type Query struct {
	Document interface{}
	Sort     string
	Limit    int64
}

func Init(collection string) *qmgo.Collection {
	database := fmt.Sprintf("%s", viper.Get("mongo.database"))
	return client.Database(database).Collection(collection)
}

func Open() {
	var (
		addr = fmt.Sprintf("%s:%d", viper.Get("mongo.addr"), viper.Get("mongo.port"))
		err  error
	)
	ctx = context.Background()
	if client, err = qmgo.NewClient(ctx, &qmgo.Config{Uri: addr}); err != nil {
		panic(err)
	}
}
