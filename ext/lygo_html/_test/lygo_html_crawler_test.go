package _test

import (
	"fmt"
	"github.com/botikasm/lygo/ext/lygo_html"
	"testing"
)

func TestBlacklist(t *testing.T) {
	urls := []string{"https://gianangelogeminiani.me",
		"https://www.facebook.com/angelo.geminiani/about?lst=1472675714%3A1472675714%3A1591518292",
		"http://facebook.com/angelo.geminiani/about?lst=1472675714%3A1472675714%3A1591518292",
	}
	for _,url:=range urls{
		match := lygo_html.UrlMatch(url, lygo_html.DefaultBlackList)
		fmt.Println(match, url)
	}
}

func TestCrawler(t *testing.T) {

	settings := new(lygo_html.HtmlCrawlerSettings)
	settings.MaxThreads = 2
	settings.StartPoints = []string{"https://gianangelogeminiani.me"}
	settings.AllowExternals = true
	settings.WhiteList = []string{"https://gianangelogeminiani.me/*"}
	settings.BlackList = []string{"https://github.com/*"}
	crawler := lygo_html.NewHtmlCrawler(settings)

	crawler.OnContent(func(content *lygo_html.HtmlCrawlerContend) {
		fmt.Println(content.Url)
		fmt.Println("\t", "error", content.Error)
		fmt.Println("\t", "links", content.Links)
		fmt.Println("\t", "blocks", len(content.Blocks))
	})

	// start and wait
	crawler.Start()
	crawler.Join()
}

func TestCrawlerLocal(t *testing.T) {

	settings := new(lygo_html.HtmlCrawlerSettings)
	settings.MaxThreads = 2
	settings.StartPoints = []string{"./pages/index.html"}

	crawler := lygo_html.NewHtmlCrawler(settings)
	crawler.Start()

}
