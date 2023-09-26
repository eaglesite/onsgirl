package main

import (
	"bytes"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/gocolly/colly/extensions"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	c := colly.NewCollector()

	c.WithTransport(&http.Transport{

		DialContext: (&net.Dialer{
			Timeout:   120 * time.Second, // 超时时间
			KeepAlive: 30 * time.Second,  // keepAlive 超时时间

		}).DialContext,
		MaxIdleConns:          500,               // 最大空闲连接数
		IdleConnTimeout:       90 * time.Second,  // 空闲连接超时
		TLSHandshakeTimeout:   120 * time.Second, // TLS 握手超时
		ExpectContinueTimeout: 10 * time.Second,
		MaxIdleConnsPerHost:   100,
		ForceAttemptHTTP2:     true,
	})

	extensions.RandomUserAgent(c)

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Add("cookie", "vtk_user_culture=cn;tempAgreementAccepted=true")
		r.Headers.Add("Accept-Language", "zh-CN,zh;q=0.9")
		r.Headers.Add("Cache-Control", "no-cache")
		r.Headers.Add("Pragma", "no-cache")
	})

	child := c.Clone()

	imagec := c.Clone()

	extensions.RandomUserAgent(child)

	extensions.RandomUserAgent(imagec)

	child.WithTransport(&http.Transport{

		DialContext: (&net.Dialer{
			Timeout:   120 * time.Second, // 超时时间
			KeepAlive: 30 * time.Second,  // keepAlive 超时时间

		}).DialContext,
		MaxIdleConns:          500,               // 最大空闲连接数
		IdleConnTimeout:       90 * time.Second,  // 空闲连接超时
		TLSHandshakeTimeout:   120 * time.Second, // TLS 握手超时
		ExpectContinueTimeout: 10 * time.Second,
		MaxIdleConnsPerHost:   100,
		ForceAttemptHTTP2:     true,
	})

	imagec.WithTransport(&http.Transport{

		DialContext: (&net.Dialer{
			Timeout:   120 * time.Second, // 超时时间
			KeepAlive: 30 * time.Second,  // keepAlive 超时时间

		}).DialContext,
		MaxIdleConns:          500,               // 最大空闲连接数
		IdleConnTimeout:       90 * time.Second,  // 空闲连接超时
		TLSHandshakeTimeout:   120 * time.Second, // TLS 握手超时
		ExpectContinueTimeout: 10 * time.Second,
		MaxIdleConnsPerHost:   100,
		ForceAttemptHTTP2:     true,
	})

	child.OnRequest(func(r *colly.Request) {
		r.Headers.Add("cookie", "vtk_user_culture=cn;tempAgreementAccepted=true")
		r.Headers.Add("Accept-Language", "zh-CN,zh;q=0.9")
		r.Headers.Add("Cache-Control", "no-cache")
		r.Headers.Add("Pragma", "no-cache")
	})

	imagec.OnRequest(func(r *colly.Request) {
		r.Headers.Add("cookie", "vtk_user_culture=cn;tempAgreementAccepted=true")
		r.Headers.Add("Accept-Language", "zh-CN,zh;q=0.9")
		r.Headers.Add("Cache-Control", "no-cache")
		r.Headers.Add("Pragma", "no-cache")
	})

	file := ""
	filem := 1
	imagec.OnResponse(func(r *colly.Response) {

		caption := strings.Split(r.Request.URL.String(), "?") // 获得刚刚#后面的信息

		u := caption[0]
		lastint := strings.LastIndex(u, "/")

		fileName := u[lastint+1:]

		f, err := os.Create(file + "/" + fileName)
		if err != nil {
			log.Println(err.Error())
		}
		io.Copy(f, bytes.NewReader(r.Body))
		fmt.Println("image-1 保存图片")
	})

	child.OnHTML("h1[class=focusbox-title]", func(el *colly.HTMLElement) {
		fmt.Println("子1-1" + "h1[class=focusbox-title]")
		file = "download/" + strconv.Itoa(filem) + "/" + strings.TrimSpace(el.Text)

		_, err := os.Stat(file)
		if err == nil {

			os.RemoveAll(file)
		}

		err = os.Mkdir(file, os.ModeType)
		if err != nil {

			log.Println(err.Error())
		}
	})

	child.OnHTML("article[class=article-content]", func(e *colly.HTMLElement) {
		if e != nil {

			//	ret := e.ChildAttrs("p > img", "src")

			e.ForEach("p", func(i int, el *colly.HTMLElement) {

				ret := el.ChildAttrs("img", "data-original")
				fmt.Println("image-1 开始")
				imagec.Visit(ret[0])
				fmt.Println("image-1 借宿")

			})

		}

	})

	// 在访问页面之前执行的回调函数
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// 在访问页面之后执行的回调函数
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL.String())
	})

	// 在访问页面时发生错误时执行的回调函数
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Error:", err)
	})

	// 在访问页面时发现图片时执行的回调函数
	c.OnHTML("a[class=thumbnail]", func(e *colly.HTMLElement) {
		url := e.Attr("href")
		if url != "" {
			fmt.Println("子1" + url)
			err := child.Visit("https://ons.ooo" + url)
			fmt.Println("子2" + url)
			if err != nil {
				log.Println("eer:" + err.Error())
			}
			log.Println("https://ons.ooo" + url)

			time.Sleep(time.Second * 10)

		}
	})

	c.OnError(func(response *colly.Response, err error) {

		log.Println(err.Error())

	})

	_, err := os.Stat("download")
	if err == nil {

		err := os.RemoveAll("download")
		if err != nil {
			return
		}
	}
	err = os.Mkdir("download", os.ModeType)
	if err != nil {
		return
	}

	for m := 4; m >= 1; m++ {

		_, err := os.Stat("download/" + strconv.Itoa(m))
		if err == nil {

			err := os.RemoveAll("download/" + strconv.Itoa(m))
			if err != nil {
				return
			}
		}
		err = os.Mkdir("download/"+strconv.Itoa(m), os.ModeType)
		if err != nil {
			return
		}
		if err != nil {
			return
		}
		t := 203
		filem = m
		/*switch m {

		case 1:
			t = 203
		case 2:
			t = 50
		case 3:
			t = 15
		case 4:
			t = 4

		}*/
		t = 2
		for i := 1; i <= t; i++ {
			fmt.Println("主1 开始")
			err := c.Visit("https://ons.ooo/type/" + strconv.Itoa(m) + "?page=" + strconv.Itoa(i))
			if err != nil {
				log.Println("err:" + err.Error())
			}

			fmt.Println("主1")
			time.Sleep(time.Minute)
		}
		time.Sleep(time.Minute)
		// 发起访问  输入你要访问的网址
	}
}
