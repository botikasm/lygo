package lygo_html

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_array"
	"github.com/botikasm/lygo/base/lygo_async"
	"github.com/botikasm/lygo/base/lygo_crypto"
	"github.com/botikasm/lygo/base/lygo_events"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_regex"
	"sync"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const (
	EventOnContent = "on_content"
)

var (
	DefaultBlackList = []string{
		"http*//*facebook.*",
		"http*//*github.*",
		"http*//*linkedin.*",
		"http*//*bitbucket.*",
		"http*//*pinterest.*",
		"http*//*instagram.*",
		"http*//*twitter.*",
		"http*//*telegram.*",
		"http*//*google.*",
		"http*//*repubblica.*",
		"http*//*akismet.*",
		"http*//*jetpack.*",
	}
)

//----------------------------------------------------------------------------------------------------------------------
//	HtmlCrawlerSettings
//----------------------------------------------------------------------------------------------------------------------

type HtmlCrawlerSettings struct {
	StartPoints             []string `json:"start_points"`
	MaxThreads              int      `json:"max_threads"`
	AllowExternals          bool     `json:"allow_externals"` // are allowed external links
	WhiteList               []string `json:"while_list"`      // always allowed
	BlackList               []string `json:"black_list"`      // never allowed
	ExcludeDefaultBlackList bool     `json:"exclude_default_black_list"`
}

func (instance *HtmlCrawlerSettings) String() string {
	return lygo_json.Stringify(instance)
}

func LoadHtmlCrawlerSettings(filename string) (*HtmlCrawlerSettings, error) {
	instance := new(HtmlCrawlerSettings)
	text, err := lygo_io.ReadTextFromFile(filename)
	if nil == err {
		lygo_json.Read(text, instance)
	} else {
		return nil, err
	}
	return instance, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	HtmlCrawlerContend
//----------------------------------------------------------------------------------------------------------------------

type HtmlCrawlerContend struct {
	Url    string           `json:"url"`
	Blocks []*SemanticBlock `json:"blocks"`
	Error  string           `json:"error"`
	Links  []string         `json:"links"`
}

//----------------------------------------------------------------------------------------------------------------------
//	HtmlCrawler
//----------------------------------------------------------------------------------------------------------------------

type HtmlCrawler struct {
	Settings *HtmlCrawlerSettings

	//-- private --//
	stopped    bool
	pool       *lygo_async.ConcurrentPool
	mux        sync.Mutex
	historyMux sync.Mutex
	history    []string
	chanURL    chan string
	chanExit   chan bool
	events     *lygo_events.Emitter
	handler    func(content *HtmlCrawlerContend)
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewHtmlCrawler(settings *HtmlCrawlerSettings) *HtmlCrawler {
	instance := new(HtmlCrawler)
	instance.stopped = true
	instance.chanURL = make(chan string)
	instance.chanExit = make(chan bool)
	instance.history = make([]string, 0)
	instance.events = lygo_events.NewEmitter()
	instance.events.On(EventOnContent, instance.onContent)
	if nil != settings {
		instance.Settings = settings
	} else {
		instance.Settings = new(HtmlCrawlerSettings)
	}
	if !instance.Settings.ExcludeDefaultBlackList {
		instance.Settings.BlackList = append(instance.Settings.BlackList, DefaultBlackList...)
	}

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *HtmlCrawler) String() string {
	if nil != instance {
		return fmt.Sprintf("Crawled %v pages", len(instance.history))
	}
	return ""
}

func (instance *HtmlCrawler) Start() {
	if nil != instance {
		instance.stopped = false
		// creates pool
		instance.pool = lygo_async.NewConcurrentPool(instance.Settings.MaxThreads)

		go instance.start()

		// add urls in settings
		for _, url := range instance.Settings.StartPoints {
			instance.chanURL <- url
		}
	}
}

func (instance *HtmlCrawler) Stop() {
	if nil != instance {
		instance.stopped = true
		if nil != instance.pool {
			instance.pool.Join()
			instance.pool = nil
		}
		instance.chanExit <- true
	}
}

func (instance *HtmlCrawler) Join() {
	if nil != instance {
		<-instance.chanExit
	}
}

func (instance *HtmlCrawler) IsWorking() bool {
	if nil != instance {
		return !instance.stopped && nil != instance.pool
	}
	return false
}

func (instance *HtmlCrawler) Crawl(path string) {
	if nil != instance {
		if instance.stopped {
			instance.Settings.StartPoints = append(instance.Settings.StartPoints, path)
		} else {
			instance.chanURL <- path
		}
	}
}

func (instance *HtmlCrawler) OnContent(callback func(event *HtmlCrawlerContend)) {
	if nil != instance && nil != callback {
		instance.handler = callback
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *HtmlCrawler) start() {
	for {
		if instance.stopped || nil == instance.pool {
			break
		}
		instance.crawl(<-instance.chanURL)
	}
}

func (instance *HtmlCrawler) onContent(event *lygo_events.Event) {
	if nil != instance.handler {
		if v, b := event.Argument(0).(*HtmlCrawlerContend); b {
			if nil != instance.handler {
				instance.handler(v)
			}
		}
	}
}

func (instance *HtmlCrawler) crawl(path string) {
	instance.mux.Lock()
	defer instance.mux.Unlock()
	if !instance.stopped && nil != instance.pool {
		startJob(0, path, instance.Settings.AllowExternals,
			instance.Settings.BlackList, instance.Settings.WhiteList,
			instance.events, instance.pool, instance.historyExists)
	}
}

func (instance *HtmlCrawler) historyExists(path string) bool {
	if nil != instance {
		instance.historyMux.Lock()
		defer instance.historyMux.Unlock()
		key := lygo_crypto.MD5(path)
		exists := lygo_array.IndexOf(key, instance.history) > -1
		if !exists {
			instance.history = append(instance.history, key)
		}
		return exists
	}
	return true
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func UrlMatch(url string, list []string) bool {
	return len(lygo_regex.WildcardIndexArray(url, list, 0)) > 0
}

func startJob(level int, path string, allowExternal bool, blackList []string, whiteList []string, events *lygo_events.Emitter, pool *lygo_async.ConcurrentPool, historyFunc func(string) bool) {
	// fmt.Println("startJob",path)
	pool.Run(func() {
		content := new(HtmlCrawlerContend)
		content.Url = path
		//get content
		parser, err := NewHtmlParser(path)
		if nil != err {
			// some error in url or network
			content.Error = err.Error()
		} else {
			// base
			rootUrl := parser.RootUrl()
			baseUrl := parser.BaseUrl()
			fullUrl := lygo_paths.Concat(baseUrl, parser.FileName())

			// url blocks
			content.Blocks = parser.SemanticBlocksAll()

			// links for children
			content.Links = parser.GetLinkURLs()
			for _, link := range content.Links {
				isExternal := len(lygo_regex.WildcardIndex(link, rootUrl+"/*", 0)) == -1
				if isExternal && !allowExternal {
					continue
				}
				isAbsolute := lygo_paths.IsAbs(link)
				if len(lygo_regex.WildcardIndexArray(link, blackList, 0)) == 0 || len(lygo_regex.WildcardIndexArray(link, whiteList, 0)) > 0 {
					// this is a good link to parse
					// fmt.Println(link, isAbsolute, isExternal)
					if isExternal && isAbsolute && !historyFunc(link) {
						go startJob(level+1, link, false, blackList, whiteList, events, pool, historyFunc)
					} else {
						if !isAbsolute {
							link = lygo_paths.Concat(baseUrl, link)
						}
						if link != fullUrl && !historyFunc(link) {
							go startJob(level, link, false, blackList, whiteList, events, pool, historyFunc)
						}
					}
				}
			} // for
		} // no error

		// raise content event
		events.EmitAsync(EventOnContent, content)
	})
}
