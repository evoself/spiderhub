package game_17173

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
	Result string `json:"result"`
	Data   []struct {
		ID               int           `json:"id"`
		ContentKey       string        `json:"contentKey"`
		ContentType      int           `json:"contentType"`
		Title            string        `json:"title"`
		TitleStyle       string        `json:"titleStyle"`
		SmallTitle       string        `json:"smallTitle"`
		ImgTitle         string        `json:"imgTitle"`
		HasImgPath       int           `json:"hasImgPath"`
		ImgPath          string        `json:"imgPath"`
		ContentFirstImg  string        `json:"contentFirstImg"`
		OriginType       int           `json:"originType"`
		Origin           string        `json:"origin"`
		OriginURL        string        `json:"originUrl"`
		Author           string        `json:"author"`
		Keywords         string        `json:"keywords"`
		KeywordsList     []string      `json:"keywordsList"`
		Description      string        `json:"description"`
		GameCodes        string        `json:"gameCodes"`
		GameCodesList    []string      `json:"gameCodesList"`
		OldGameCodes     string        `json:"oldGameCodes"`
		OldGameCodesList []int         `json:"oldGameCodesList"`
		Tags             string        `json:"tags"`
		TagsList         []string      `json:"tagsList"`
		Weight           int           `json:"weight"`
		WeightSet        int           `json:"weightSet"`
		Status           int           `json:"status"`
		CopyrightInfo    string        `json:"copyrightInfo"`
		PublishTime      int64         `json:"publishTime"`
		UpdateTime       int64         `json:"updateTime"`
		CreateTime       int64         `json:"createTime"`
		CreateUserID     int           `json:"createUserId"`
		HitCount         interface{}   `json:"hitCount"`
		SupportCount     interface{}   `json:"supportCount"`
		OpposeCount      interface{}   `json:"opposeCount"`
		UpdateUserID     int           `json:"updateUserId"`
		ReleaseUserID    int           `json:"releaseUserId"`
		Plugin           interface{}   `json:"plugin"`
		RelateList       interface{}   `json:"relateList"`
		ShiftTitle       interface{}   `json:"shiftTitle"`
		PageURL          string        `json:"pageUrl"`
		CategoryID       int           `json:"categoryId"`
		DefaultCategory  int           `json:"defaultCategory"`
		CategoryIds      string        `json:"categoryIds"`
		CategoryIdsList  []int         `json:"categoryIdsList"`
		CategoryName     string        `json:"categoryName"`
		CategoryURL      string        `json:"categoryUrl"`
		DomainID         int           `json:"domainId"`
		ChannelID        int           `json:"channelId"`
		ChannelCode      int           `json:"channelCode"`
		ChannelName      interface{}   `json:"channelName"`
		ChannelURL       interface{}   `json:"channelUrl"`
		ContentText      string        `json:"contentText"`
		ImageGroupValue  interface{}   `json:"imageGroupValue"`
		Content          string        `json:"content"`
		HasComment       int           `json:"hasComment"`
		HasImage         int           `json:"hasImage"`
		HasVedio         int           `json:"hasVedio"`
		ImageList        []string      `json:"imageList"`
		VideoList        []interface{} `json:"videoList"`
		Extension1       interface{}   `json:"extension1"`
		Extension2       interface{}   `json:"extension2"`
		Extension3       interface{}   `json:"extension3"`
		Extension4       interface{}   `json:"extension4"`
		Extension5       interface{}   `json:"extension5"`
		GameIDList       []string      `json:"gameIdList"`
	} `json:"data"`
	Status     int `json:"status"`
	TotalCount int `json:"totalCount"`
}

var running = true

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
	req.Header.Set("Host", "news.17173.com")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Referer", "http://www.17173.com")
	req.Header.Set("Accept-Encoding", "*")
	req.Header.Set("Accept-Language", "zh-CN,zh;q=0.9")
	req.Header.Set("User-Agent", ua.UserAgentMobile())
	resp, err = c.Do(req)
	if err != nil {
		models.NewLog(err.Error()).Save()
	}
	if resp == nil {
		return
	}
	defer resp.Body.Close()
	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	_ = json.Unmarshal(result, &m)

	// 获取数据失败
	if m.Result == "failure" {
		running = false
		return
	}
	for _, i := range m.Data {
		p := models.NewPage()
		p.Url = i.PageURL
		p.Title = i.Title
		p.Source = "17173"
		p.Image = "http:" + i.ImgPath
		p.Category = "游戏"
		p.Date, _ = strconv.Atoi(time.Now().Format("20060102"))
		p.ExtractTime = gotime.FormatDatetime(time.Now(), gotime.TT)
		p.PublishTime = gotime.FormatDatetime(time.Unix(i.PublishTime, 0), gotime.TT)
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
	i := 1
	for running {
		if running == false {
			return
		}
		target := "http://news.17173.com/data/content/list.json?pageSize=10&pageNo=" + strconv.Itoa(i)
		extract(target)
		i++
	}
}
