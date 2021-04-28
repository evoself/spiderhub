package guokr

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"spiderhub/models"
	"strconv"
	"time"
)

var running = true

type model struct {
	Author struct {
		Avatar struct {
			Large  string `json:"large"`
			Normal string `json:"normal"`
			Small  string `json:"small"`
		} `json:"avatar"`
		Nickname string `json:"nickname"`
		Title    string `json:"title"`
		Ukey     string `json:"ukey"`
	} `json:"author"`
	Authors []struct {
		Avatar struct {
			Large  string `json:"large"`
			Normal string `json:"normal"`
			Small  string `json:"small"`
		} `json:"avatar"`
		Nickname string `json:"nickname"`
		Title    string `json:"title"`
		Ukey     string `json:"ukey"`
	} `json:"authors"`
	BaiduEditor struct {
	} `json:"baidu_editor"`
	Category struct {
	} `json:"category"`
	Channels          []interface{} `json:"channels"`
	DateCreated       time.Time     `json:"date_created"`
	DateModified      time.Time     `json:"date_modified"`
	DatePublished     time.Time     `json:"date_published"`
	ID                int           `json:"id"`
	Image             string        `json:"image"`
	ImagePortrait     string        `json:"image_portrait"`
	IsEditorRecommend bool          `json:"is_editor_recommend"`
	IsLiyanArticle    bool          `json:"is_liyan_article"`
	IsPublished       bool          `json:"is_published"`
	IsReplyable       bool          `json:"is_replyable"`
	MinisiteKey       interface{}   `json:"minisite_key"`
	RepliesCount      int           `json:"replies_count"`
	SmallImage        string        `json:"small_image"`
	Subject           struct {
		ArticlesCount int       `json:"articles_count"`
		DateCreated   time.Time `json:"date_created"`
		Key           string    `json:"key"`
		MinisiteKey   string    `json:"minisite_key"`
		Name          string    `json:"name"`
		SortScore     int       `json:"sort_score"`
		SubjectType   string    `json:"subject_type"`
	} `json:"subject"`
	SubjectKey   string `json:"subject_key"`
	Summary      string `json:"summary"`
	Title        string `json:"title"`
	UkeyAuthor   string `json:"ukey_author"`
	VideoContent string `json:"video_content"`
}

func extract(target string) {
	resp, err := http.Get(target)
	if err != nil {
		_, _ = models.NewLog(err.Error()).Save()
	}
	if resp == nil {
		return
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	var m []model
	_ = json.Unmarshal(result, &m)
	// 获取数据失败
	if len(m) == 0 {
		running = false
	}
	for _, i := range m {
		fmt.Println(i)
	}
}

// 娱乐频道：
//https://interface.sina.cn/ent/feed.d.json?ch=ent&col=ent&act=more&t=1484477669001&show_num=10&page=4
//参数说明：
//ch:频道
//娱乐：ent
//体育：sports
//科技：tech
//教育：edu
//健康：health
//时尚：fashion
//博客：blog
//col：分类
//show_num
//page
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

	i := 1
	for running {
		target := "https://www.guokr.com/beta/proxy/science_api/articles?limit=30&page=" + strconv.Itoa(i)
		extract(target)
		i++
	}
}