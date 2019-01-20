package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"github.com/gocolly/colly"
	"github.com/robertkrimen/otto"
)



func main() {
	c := colly.NewCollector(
		colly.AllowedDomains("xiai123.com","ting.xiai123.com"),
	)

	titleMap := map[string]string{};

	runtime := otto.New()
	f:=`function SMusic(a) {
		return a.musicList
	};`

	c.OnHTML("title", func(e *colly.HTMLElement) {
		title := e.Text
		url := e.Response.Request.URL.Path
		ext := filepath.Ext(url)
		id := strings.TrimSuffix(filepath.Base(url),ext)
		titleMap[id]=title;
	})

	c.OnHTML("script[src]", func(e *colly.HTMLElement) {
		src := e.Attr("src")
		if strings.HasPrefix(src,"http://ting.xiai123.com/") {
			c.Visit(src)
		}
	})

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		c.Visit(e.Request.AbsoluteURL(link))
	})

	c.OnResponse(func(resp *colly.Response){
		if resp.Headers.Get("Content-Type") == "application/javascript"{
			bytes := resp.Body

			if _, err := runtime.Run(f); err != nil {
				panic(err)
			}

			jsvalue, err := runtime.Run(string(bytes))
			if err != nil {
				panic(err)
			}

			export, err := jsvalue.Export()
			if err != nil {
				panic(err)
			}

			v, _ := export.([]map[string]interface {})
			for _,value :=range v {
				s := value["src"].(string)
				if strings.Index(s,"kewen") > -1 {
					c.Visit(s)
				}
			}

		}

		if resp.Headers.Get("Content-Type") == "audio/mpeg"{
			url,_ := filepath.Rel("/mp3/",resp.Request.URL.RawPath)

			id := filepath.Base(filepath.Dir(url))
			title := titleMap[id]

			if title != ""{
				url=strings.Replace(url,id,title,-1)
			}

			os.MkdirAll("download/"+filepath.Dir(url),os.ModePerm)

			go resp.Save("download/"+url)
		}
	})


	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("authority","www.jianshu.com")
		r.Headers.Set("pragma","no-cache")
		r.Headers.Set("cache-control","no-cache")
		r.Headers.Set("upgrade-insecure-requests","1")
		r.Headers.Set("user-agent","Mozilla/5.0(Macintosh;IntelMacOSX10_13_6)AppleWebKit/537.36(KHTML,likeGecko)Chrome/71.0.3578.98Safari/537.36")
		r.Headers.Set("accept","text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8")
		r.Headers.Set("accept-encoding","gzip,deflate,br")
		r.Headers.Set("accept-language","en,zh-CN;q=0.9,zh;q=0.8")
		r.Headers.Set("cookie","read_mode=day;default_font=font2;_ga=GA1.2.843094157.1512263796;locale=zh-CN;__yadk_uid=3Cax2lRRsIv7GrEplxjuk0mtvoAB3Dkc;signin_redirect=https%3A%2F%2Fwww.jianshu.com%2Fp%2Fa2ae558d0e41;_m7e_session_core=0b2611f45826123683949cff8487ced2;Hm_lvt_0c0e9d9b1e7d617b3e6842e85b9fb068=1547870930,1547904764,1547910458,1547965147;Hm_lpvt_0c0e9d9b1e7d617b3e6842e85b9fb068=1547965182;sensorsdata2015jssdkcross=%7B%22distinct_id%22%3A%2216019f2465672d-0b65fce48ccbb8-173d6d56-1049088-16019f2465720e%22%2C%22%24device_id%22%3A%2216019f2465672d-0b65fce48ccbb8-173d6d56-1049088-16019f2465720e%22%2C%22props%22%3A%7B%22%24latest_traffic_source_type%22%3A%22%E7%9B%B4%E6%8E%A5%E6%B5%81%E9%87%8F%22%2C%22%24latest_referrer%22%3A%22%22%2C%22%24latest_referrer_host%22%3A%22%22%2C%22%24latest_search_keyword%22%3A%22%E6%9C%AA%E5%8F%96%E5%88%B0%E5%80%BC_%E7%9B%B4%E6%8E%A5%E6%89%93%E5%BC%80%22%2C%22%24latest_utm_campaign%22%3A%22maleskine%22%2C%22%24latest_utm_source%22%3A%22recommendation%22%2C%22%24latest_utm_medium%22%3A%22pc_all_hots%22%2C%22%24latest_utm_content%22%3A%22note%22%7D%2C%22first_id%22%3A%22%22%7D")
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("http://xiai123.com/2018qiu-rjyw3s.html")
}