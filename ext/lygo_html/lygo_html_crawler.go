package lygo_html

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_array"
	"github.com/botikasm/lygo/base/lygo_async"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_regex"
	"sync"
)

//----------------------------------------------------------------------------------------------------------------------
//	HtmlCrawlerSettings
//----------------------------------------------------------------------------------------------------------------------

type HtmlCrawlerSettings struct {
	StartPoints    []string `json:"start_points"`
	MaxThreads     int      `json:"max_threads"`
	AllowExternals bool     `json:"allow_externals"` // are allowed external links
	WhiteList      []string `json:"while_list"`      // always allowed
	BlackList      []string `json:"black_list"`      // never allowed
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
	if nil != settings {
		instance.Settings = settings
	} else {
		instance.Settings = new(HtmlCrawlerSettings)
	}

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *HtmlCrawler) String() string {
	if nil != instance {

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

func (instance *HtmlCrawler) crawl(path string) {
	instance.mux.Lock()
	defer instance.mux.Unlock()
	if !instance.stopped && nil != instance.pool {
		startJob(0, path, instance.Settings.AllowExternals,
			instance.Settings.BlackList, instance.Settings.WhiteList,
			instance.pool, instance.historyExists)
	}
}

func (instance *HtmlCrawler) historyExists(path string) bool {
	if nil != instance {
		instance.historyMux.Lock()
		defer instance.historyMux.Unlock()
		exists := lygo_array.IndexOf(path, instance.history) > -1
		if !exists {
			instance.history = append(instance.history, path)
		}
		return exists
	}
	return true
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func startJob(level int, path string, allowExternal bool, blackList []string, whiteList []string, pool *lygo_async.ConcurrentPool, historyFunc func(string) bool) {
	pool.Run(func() {
		//get content
		parser, err := NewHtmlParser(path)
		if nil != err {
			// some error in url or network
			fmt.Println("ERROR", err, path)
		} else {
			// base
			baseUrl := parser.BaseUrl()
			fullUrl := lygo_paths.Concat(baseUrl, parser.FileName())

			// url blocks
			blocks := parser.SemanticBlocksAll()
			if len(blocks) > 0 {

			}

			// links for children
			links := parser.GelLinkURLs()
			for _, link := range links {
				isExternal := len(lygo_regex.WildcardIndex(link, baseUrl+"*", 0)) == 0
				if isExternal && !allowExternal {
					continue
				}
				isAbsolute := lygo_paths.IsAbs(link)
				if len(lygo_regex.WildcardIndexArray(link, blackList, 0)) == 0 || len(lygo_regex.WildcardIndexArray(link, whiteList, 0)) > 0 {
					// this is a good link to parse
					fmt.Println(link, isAbsolute, isExternal)
					if isExternal && isAbsolute && !historyFunc(link) {
						startJob(level+1, link, false, blackList, whiteList, pool, historyFunc)
					} else {
						if !isAbsolute {
							link = lygo_paths.Concat(baseUrl, link)
						}
						if link != fullUrl && !historyFunc(link) {
							startJob(level, link, false, blackList, whiteList, pool, historyFunc)
						}
					}
				}
			}
		}
	})
}
