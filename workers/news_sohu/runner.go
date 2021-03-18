package news_sohu

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
	"time"
)

const (
	host   = "https://m.sohu.com/"
	target = "https://v2.sohu.com/integration-api/mix/region/101?page=1&size=2"
)

type model struct {
	Data []struct {
		Brief         string `json:"brief,omitempty"`
		ImageInfoList []struct {
			Width  int    `json:"width"`
			URL    string `json:"url"`
			Height int    `json:"height"`
		} `json:"imageInfoList"`
		Images             []string `json:"images"`
		CmsID              int      `json:"cmsId"`
		MobileTitle        string   `json:"mobileTitle"`
		MobilePersonalPage string   `json:"mobilePersonalPage"`
		Type               int      `json:"type"`
		AuthorID           int      `json:"authorId"`
		AuthorPic          string   `json:"authorPic"`
		Title              string   `json:"title"`
		URL                string   `json:"url"`
		Cover              string   `json:"cover,omitempty"`
		PublicTime         int64    `json:"publicTime"`
		AuthorName         string   `json:"authorName"`
		ID                 int      `json:"id"`
		Scm                string   `json:"scm"`
		PersonalPage       string   `json:"personalPage"`
		BigCover           string   `json:"bigCover,omitempty"`
		ResourceType       int      `json:"resourceType"`
		VideoInfo          struct {
			Duration      int    `json:"duration"`
			SmartDuration string `json:"smartDuration"`
			Site          int    `json:"site"`
			SofaInfo      []struct {
				VideoWidth   int    `json:"videoWidth"`
				PlayID       int    `json:"playId"`
				VideoURL     string `json:"videoUrl"`
				Rate         string `json:"rate"`
				VideoLevelID int    `json:"videoLevelId"`
				VideoSize    int    `json:"videoSize"`
				X265         bool   `json:"x265"`
				VideoHeight  int    `json:"videoHeight"`
				VideoStatus  int    `json:"videoStatus"`
			} `json:"sofaInfo"`
			Width   int `json:"width"`
			VideoID int `json:"videoId"`
			Height  int `json:"height"`
		} `json:"videoInfo,omitempty"`
		TagList []struct {
			Name   string `json:"name"`
			ID     int    `json:"id"`
			Type   int    `json:"type"`
			Status int    `json:"status"`
		} `json:"tagList,omitempty"`
		ImageNum int `json:"imageNum,omitempty"`
	} `json:"data"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}

// 开始提取数据
func extract(target string) {
	var (
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

	// 获取数据失败
	if m.Message != "success" {
		return
	}
	for _, i := range m.Data {
		p := models.NewPage()
		p.Url = host + i.URL
		p.Title = i.Title
		p.Source = "搜狐"
		p.ExtractTime = gotime.FormatDatetime(time.Now(), gotime.TT)
		p.PublishTime = gotime.FormatDatetime(time.Unix(i.PublicTime, 0), gotime.TT)
		p.Hash = hash.Sha256String(p.Url)
		p.Id = time.Now().UnixNano() + int64(rand.RandInt(100, 999))
		res, _ := p.Save()
		if res != nil {
			color.Green.Println(p.Source + "-" + p.Title)
		} else {
			color.Red.Println(p.Source + "-" + p.Title)
		}
	}
}

func Run() {
	categoryMap := make(map[string]string)
	//microsecond := fmt.Sprintf("%d", time.Now().UnixNano()/1000)
	categoryMap["tech"] = "科技"
	categoryMap["ent"] = "娱乐"
	categoryMap["sports"] = "体育"
	categoryMap["edu"] = "教育"
	categoryMap["health"] = "健康"
	categoryMap["fashion"] = "时尚"
	categoryMap["blog"] = "博客"

	//for k, v := range categoryMap {
	//
	//}
}
