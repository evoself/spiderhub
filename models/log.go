package models

import (
	"context"
	"github.com/qiniu/qmgo"
	"path/filepath"
	"runtime"
	"spiderhub/pkg/mgo"
	"strings"
	"time"
)

const logCollection = "log"

type log struct {
	Code     int    `bson:"code"`
	Message  string `bson:"message"`
	Filename string `bson:"filename"`
	Funcname string `bson:"funcname"`
	Line     int    `bson:"line"`
	Date     string `bson:"date"`
	Created  string `bson:"created"`
}

// 告警方法
func NewLog(msg string) *log {
	// 当前时间
	currentTime := time.Now()

	// 定义 文件名、行号、方法名
	fileName, line, functionName := "?", 0, "?"

	pc, fileName, line, ok := runtime.Caller(2)
	if ok {
		functionName = runtime.FuncForPC(pc).Name()
		functionName = filepath.Ext(functionName)
		functionName = strings.TrimPrefix(functionName, ".")
	}

	newlog := &log{
		Message:  msg,
		Filename: fileName,
		Line:     line,
		Funcname: functionName,
		Date:     currentTime.Format("20060102"),
		Created:  currentTime.Format("2006-01-02 15:04:05"),
	}
	return newlog
}

func (l *log) Save() (*qmgo.InsertOneResult, error) {
	l.Date = time.Now().Format("20060102")
	l.Created = time.Now().Format("2006-01-02 15:04:05")
	return mgo.Init(logCollection).InsertOne(context.Background(), l)
}
