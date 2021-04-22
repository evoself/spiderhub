package news_36kr

import (
	"bytes"
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
	Code int `json:"code"`
	Data struct {
		ItemList []struct {
			ItemID           int64 `json:"itemId"`
			ItemType         int   `json:"itemType"`
			TemplateMaterial struct {
				ItemID       int64  `json:"itemId"`
				TemplateType int    `json:"templateType"`
				WidgetImage  string `json:"widgetImage"`
				PublishTime  int64  `json:"publishTime"`
				WidgetTitle  string `json:"widgetTitle"`
				AuthorName   string `json:"authorName"`
				NavName      string `json:"navName"`
			} `json:"templateMaterial"`
			Route string `json:"route"`
		} `json:"itemList"`
		PageCallback string `json:"pageCallback"`
		HasNextPage  int    `json:"hasNextPage"`
	} `json:"data"`
}

var (
	b            = true
	running      = true
	pageCallback string
)

// 开始提取数据
func extract(target string) {
	var (
		m      = &model{}
		c      = &http.Client{Timeout: time.Second * 5}
		result []byte
		resp   *http.Response
		err    error
		param  = make(map[string]interface{})
		params = make(map[string]interface{})
	)
	if b {
		param["pageCallback"] = "eyJmaXJzdElkIjozMTQ1NTkzLCJsYXN0SWQiOjMxNDU1MTYsImZpcnN0Q3JlYXRlVGltZSI6MTYwODI4Nzc2MTkwMiwibGFzdENyZWF0ZVRpbWUiOjE2MDgyODY0OTc1NDh9"
	} else {
		param["pageCallback"] = pageCallback
	}
	param["pageSize"] = 10
	param["pageEvent"] = 1
	param["platformId"] = 2
	param["siteId"] = 1
	params["partner_id"] = "wap"
	params["timestamp"] = time.Now().UnixNano() / 1000
	params["param"] = param
	byteArray, _ := json.Marshal(params)
	req, _ := http.NewRequest("POST", target, bytes.NewBuffer(byteArray))
	req.Header.Set("Cookie", "Hm_lvt_1684191ccae0314c6254306a8333d090=1608255147,1608255183,1608255193,1608256180; Hm_lvt_713123c60a0e86982326bae1a51083e1=1608255147,1608255183,1608255193,1608256180; sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%22175cefcf9654c1-0d5017d88cabb5-445e6c-1296000-175cefcf96654d%22%2C%22%24device_id%22%3A%22175cefcf9654c1-0d5017d88cabb5-445e6c-1296000-175cefcf96654d%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_referrer%22%3A%22%22%2C%22%24latest_referrer_host%22%3A%22%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%7D%7D; _ss_pp_id=1e18a005e829d6cce771605466691960; _td=ecf117ee-9498-4e83-8f4d-3ee813147cae; Hm_lpvt_1684191ccae0314c6254306a8333d090=1608300677; Hm_lpvt_713123c60a0e86982326bae1a51083e1=1608300677; acw_tc=2760825b16083017175743374ec7c10744408420e2d116a0d8477b4b3d8206")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Host", "gateway.36kr.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Length", "260")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
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
	if len(m.Data.ItemList) == 0 {
		running = false
		return
	}
	for _, i := range m.Data.ItemList {
		p := models.NewPage()
		p.Url = "https://36kr.com/p/" + strconv.FormatInt(i.ItemID, 10)
		p.Title = i.TemplateMaterial.WidgetTitle
		p.Image = i.TemplateMaterial.WidgetImage
		p.Source = "36氪"
		p.Category = "科技"
		p.Date, _ = strconv.Atoi(time.Now().Format("20060102"))
		p.ExtractTime = gotime.FormatDatetime(time.Now(), gotime.TT)
		tm := time.Unix(i.TemplateMaterial.PublishTime/1000, 0)
		p.PublishTime = gotime.FormatDatetime(tm, gotime.TT)
		p.Hash = hash.Sha256String(p.Url)
		p.Id = time.Now().UnixNano() + int64(rand.RandInt(100, 999))
		res, _ := p.Save()
		if res != nil {
			color.Green.Println(p.Source + "-" + p.Title)
		} else {
			color.Red.Println(p.Source + "-" + p.Title)
		}
	}
	pageCallback = m.Data.PageCallback
	b = false
}

func Run() {
	target := "https://gateway.36kr.com/api/mis/nav/home/flow/forWap"
	for running {
		if running == false {
			return
		}
		extract(target)
	}
}
