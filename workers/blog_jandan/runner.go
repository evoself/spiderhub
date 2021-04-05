package blog_jandan

import (
	"github.com/gocolly/colly"
	"github.com/gohp/goutils/color"
	"github.com/gohp/goutils/gotime"
	"github.com/gohp/goutils/hash"
	"github.com/gohp/goutils/rand"
	mrand "math/rand"
	"spiderhub/models"
	"spiderhub/pkg/ua"
	"strconv"
	"strings"
	"time"
)

func Run() {
	var (
		p        = models.NewPage()
		n        = 1
		running  = true
		callback = func(target string) {
			c := colly.NewCollector()
			c.OnRequest(func(r *colly.Request) {
				r.Headers.Set("Host", "i.jandan.net")
				r.Headers.Set("Connection", "keep-alive")
				r.Headers.Set("Accept", "*/*")
				r.Headers.Set("Accept-Encoding", "gzip, deflate")
				r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
				r.Headers.Set("Upgrade-Insecure-Requests", "1")
				r.Headers.Set("User-Agent", ua.UserAgentMobile())
			})
			c.OnHTML(".postinfo", func(e *colly.HTMLElement) {
				arr := strings.Split(e.Text, "@")
				datetime := strings.Trim(arr[1], " ")
				if strings.Index(datetime, "下午") > 0 {
					datetime = strings.Replace(datetime, "下午", "PM", -1)
				} else {
					datetime = strings.Replace(datetime, "上午", "AM", -1)
				}

				timeLayout := "2006.01.02 , 03:04 PM"
				parseTime, _ := time.Parse(timeLayout, datetime)
				//将时间戳设置成种子数
				mrand.Seed(time.Now().UnixNano())
				//生成10个0-99之间的随机数
				p.PublishTime = parseTime.Format("2006-01-02 15:04") + ":" + strconv.Itoa(mrand.Intn(50)+10)
			})
			_ = c.Visit(target)
			c.Wait()
		}
	)
	for running {
		target := "http://i.jandan.net/page/" + strconv.Itoa(n)
		c := colly.NewCollector()
		c.OnRequest(func(r *colly.Request) {
			r.Headers.Set("Host", "i.jandan.net")
			r.Headers.Set("Connection", "keep-alive")
			r.Headers.Set("Accept", "*/*")
			r.Headers.Set("Accept-Encoding", "gzip, deflate")
			r.Headers.Set("Accept-Language", "zh-CN,zh;q=0.9")
			r.Headers.Set("User-Agent", ua.UserAgentMobile())
		})
		// 图片
		c.OnHTML(".post .posthit > .thumb_s > a", func(e *colly.HTMLElement) {
			img, _ := e.DOM.Find("img").Attr("data-original")
			array := strings.Split(img, "!")
			if len(array) == 2 {
				p.Image = "https:" + array[0]
			} else {
				p.Image = array[0]
			}
		})
		c.OnHTML(".post .posthit h2 > a", func(e *colly.HTMLElement) {
			p.Url = e.Attr("href")
			p.Title = e.Text
			p.Source = "煎蛋"
			p.Category = "社区"
			p.Date, _ = strconv.Atoi(time.Now().Format("20060102"))
			p.ExtractTime = gotime.FormatDatetime(time.Now(), gotime.TT)
			p.Hash = hash.Sha256String(p.Url)
			p.Id = time.Now().UnixNano() + int64(rand.RandInt(100, 999))
			callback(p.Url)
			res, _ := p.Save()
			if res != nil {
				color.Green.Println(p.Source + "-" + p.Title)
			} else {
				color.Red.Println(p.Source + "-" + p.Title)
			}
		})
		// 数据完
		c.OnHTML(".postnohit", func(e *colly.HTMLElement) {
			running = false
		})
		c.Visit(target)
		c.Wait()
		n++
	}
}
