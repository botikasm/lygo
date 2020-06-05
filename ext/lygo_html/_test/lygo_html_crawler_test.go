package _test

import (
	"github.com/botikasm/lygo/ext/lygo_html"
	"testing"
)

func TestCrawler(t *testing.T) {

	settings := new(lygo_html.HtmlCrawlerSettings)
	settings.MaxThreads = 2
	settings.StartPoints = []string{"https://gianangelogeminiani.me"}
	settings.AllowExternals = true
	settings.WhiteList = []string{"https://gianangelogeminiani.me/*"}
	settings.BlackList = []string{"https://github.com/*"}
	crawler := lygo_html.NewHtmlCrawler(settings)

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
