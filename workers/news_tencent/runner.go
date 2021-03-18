package news_tencent

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

type model struct {
	Ret  int    `json:"ret"`
	Msg  string `json:"msg"`
	Data struct {
		List []struct {
			CmsID           string   `json:"cms_id" bson:"cms_id"`
			Title           string   `json:"title" bson:"title"`
			Subtitle        string   `json:"subtitle" bson:"subtitle"`
			Url             string   `json:"url" bson:"url"`
			ThumbNail       string   `json:"thumb_nail" bson:"thumbnail"`
			ThumbNail2X     string   `json:"thumb_nail_2x" bson:"thumbnail_2x"`
			TopBigImg       []string `json:"top_big_img" bson:"top_big_img"`
			CategoryID      string   `json:"category_id" bson:"category_id"`
			CategoryName    string   `json:"category_name" bson:"category_name"`
			CategoryCn      string   `json:"category_cn" bson:"category_cn"`
			SubCategoryID   string   `json:"sub_category_id" bson:"sub_category_id"`
			SubCategoryName string   `json:"sub_category_name" bson:"sub_category_name"`
			SubCategoryCn   string   `json:"sub_category_cn" bson:"sub_category_cn"`
			Status          int      `json:"status" bson:"status"`
			Tags            []struct {
				TagID    string `json:"tag_id" bson:"tag_id"`
				TagWord  string `json:"tag_word" bson:"tag_word"`
				TagScore string `json:"tag_score" bson:"tag_score"`
			} `json:"tags" bson:"tags"`
			MediaID       string `json:"media_id" bson:"media_id"`
			MediaName     string `json:"media_name" bson:"media_name"`
			Point         string `json:"point" bson:"point"`
			ArticleType   int    `json:"article_type" bson:"article_type"`
			PoolName      string `json:"pool_name" bson:"pool_name"`
			SecurityField int    `json:"security_field" bson:"security_field"`
			ArticleID     string `json:"article_id" bson:"article_id"`
			Source        string `json:"source" bson:"source"`
			CommentID     string `json:"comment_id" bson:"comment_id"`
			CommentNum    string `json:"comment_num" bson:"comment_num"`
			CreateTime    string `json:"create_time" bson:"create_time"`
			UpdateTime    string `json:"update_time" bson:"update_time"`
			PublishTime   string `json:"publish_time" bson:"publish_time"`
			ImgExpType    string `json:"img_exp_type" bson:"img_exp_type"`
			Img           string `json:"img" bson:"img"`
			UrlHash       string `bson:"url_hash"`
		} `json:"list"`
	} `json:"data"`
}

var running = true

// 开始提取数据
func extract(target, category string) {
	var (
		p      = models.NewPage()
		m      = &model{}
		c      = &http.Client{Timeout: time.Second * 5}
		result []byte
		resp   *http.Response
		err    error
	)

	req, _ := http.NewRequest("GET", target, nil)
	req.Header.Set("Host", "i.news.qq.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Cache-Control", "no-cache")
	req.Header.Set("Referer", target)
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
		_, _ = models.NewLog(err.Error()).Save()
	}
	_ = json.Unmarshal(result, &m)

	if len(m.Data.List) == 0 {
		running = false
		return
	}
	for _, i := range m.Data.List {
		// https://kuaibao.qq.com/a/20201215A0A0UT00
		p.Url = "https://view.inews.qq.com/a/" + i.CmsID
		//p.Url = "https://view.inews.qq.com/s/" + i.CmsID
		p.Title = i.Title
		p.Source = "腾讯"
		p.Image = i.Img
		p.Category = category
		p.Date, _ = strconv.Atoi(time.Now().Format("20060102"))
		p.ExtractTime = gotime.FormatDatetime(time.Now(), gotime.TT)
		p.PublishTime = i.PublishTime
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
	categoryMap["tech"] = "科技"
	categoryMap["games"] = "游戏"
	categoryMap["fashion"] = "时尚"
	categoryMap["history"] = "历史"
	categoryMap["finance_stock"] = "股票"
	categoryMap["digi"] = "数码"
	categoryMap["auto"] = "汽车"
	categoryMap["sports"] = "体育"
	categoryMap["cul"] = "文化"
	categoryMap["nba"] = "体育"
	categoryMap["ent"] = "娱乐"
	categoryMap["finance"] = "财经"
	categoryMap["edu"] = "教育"
	categoryMap["world"] = "国际"
	categoryMap["milite"] = "军事"
	categoryMap["astro"] = "星座"

	for k, v := range categoryMap {
		i := 1
		if running == false {
			running = true
			continue
		}
		for running {
			target := "https://i.news.qq.com/trpc.qqnews_web.kv_srv.kv_srv_http_proxy/list?" +
				"srv_id=pc&sub_srv_id=" + k +
				"&strategy=1&ext={\"pool\":[\"high\",\"top\"],\"is_filter\":10,\"check_type\":true}" +
				"&offset=" + strconv.Itoa(i) + "&limit=199"
			extract(target, v)
			i++
		}
	}
}
