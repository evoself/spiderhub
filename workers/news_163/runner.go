package news_163

// 本地新闻https://3g.163.com/touch/jsonp/article/local/%E6%B9%96%E5%8D%97/10-50.html
// 家居 http://home.3g.163.com/home/mobile/interface/CMSNews/00108835/1.html
import (
	"encoding/json"
	"github.com/gohp/goutils/color"
	"github.com/gohp/goutils/gotime"
	"github.com/gohp/goutils/hash"
	"github.com/gohp/goutils/rand"
	"io/ioutil"
	"net/http"
	"spiderhub/models"
	"spiderhub/pkg/ua"
	"strconv"
	"time"
)

type Model struct {
	//BA8D4A3Rwangning []struct {
	Data []struct {
		LiveInfo     interface{} `json:"liveInfo" bson:"live_info"`
		Docid        string      `json:"docid" bson:"doc_id"`
		Source       string      `json:"source" bson:"source"`
		Title        string      `json:"title" bson:"title"`
		Priority     int         `json:"priority" bson:"priority"`
		HasImg       int         `json:"hasImg" bson:"has_img"`
		URL          string      `json:"url" bson:"url"`
		CommentCount int         `json:"commentCount" bson:"comment_Count"`
		Imgsrc3Gtype string      `json:"imgsrc3gtype" bson:"img_src_3g_type"`
		Stitle       string      `json:"stitle" bson:"stitle"`
		Digest       string      `json:"digest" bson:"digest"`
		Imgsrc       string      `json:"imgsrc" bson:"img_src"`
		Ptime        string      `json:"ptime" bson:"ptime"`
	} `json:"BA8D4A3Rwangning"`
}

// 开始提取数据
func extract(target, category string) {
	var (
		m      = &Model{}
		c      = &http.Client{Timeout: time.Second * 5}
		result []byte
		resp   *http.Response
		err    error
	)
	req, _ := http.NewRequest("GET", target, nil)
	req.Header.Set("Host", "3g.163.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Referer", "https://3g.163.com")
	req.Header.Set("Accept-Encoding", "*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("User-Agent", ua.UserAgentMobile())
	resp, err = c.Do(req)
	if err != nil {
		_, _ = models.NewLog(err.Error()).Save()
	}
	if resp == nil {
		return
	}
	defer resp.Body.Close()
	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	length := len(string(result))
	// 从第9位开始截取
	arrayByte := result[9:(length - 1)]
	_ = json.Unmarshal(arrayByte, &m)
	for _, i := range m.Data {
		p := models.NewPage()
		p.Url = i.URL
		p.Title = i.Title
		p.Source = "网易"
		p.Image = i.Imgsrc
		p.Category = category
		p.Hash = hash.Sha256String(p.Url)
		p.Date, _ = strconv.Atoi(time.Now().Format("20060102"))
		p.ExtractTime = gotime.FormatDatetime(time.Now(), gotime.TT)
		p.PublishTime = i.Ptime
		p.Hash = hash.Sha256String(p.Url)
		p.Id = time.Now().UnixNano() + int64(rand.RandInt(100, 999))
		res, _ := p.Save()
		if res != nil {
			color.Green.Println(p.Source + "-" + p.Title)
		} else {
			color.Red.Println(p.Source + "-" + p.Title)
		}
		time.Sleep(time.Second * 1)
	}
}

func Run() {
	categoryMap := make(map[string]string)
	categoryMap["BA8D4A3Rwangning"] = "科技"
	categoryMap["BA10TA81wangning"] = "娱乐"
	categoryMap["BA8E6OEOwangning"] = "体育"
	categoryMap["BA8EE5GMwangning"] = "财经"
	categoryMap["BA8DOPCSwangning"] = "汽车"
	categoryMap["BAI67OGGwangning"] = "军事"
	categoryMap["BAI6JOD9wangning"] = "数码"
	categoryMap["BAI6RHDKwangning"] = "游戏"
	categoryMap["BA8FF5PRwangning"] = "教育"
	categoryMap["BA8F6ICNwangning"] = "时尚"
	categoryMap["BDC4QSV3wangning"] = "健康"

	total := 300
	for k, v := range categoryMap {
		for i := 0; i <= total; i++ {
			target := "https://3g.163.com/touch/reconstruct/article/list/" + k + "/" +
				strconv.FormatInt(int64(i), 10) + "-20.html"
			extract(target, v)
		}
	}
}
