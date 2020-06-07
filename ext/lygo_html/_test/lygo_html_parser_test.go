package _test

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/ext/lygo_html"
	"golang.org/x/net/html"
	"testing"
)

func TestCreate(t *testing.T) {

	parser, err := lygo_html.NewHtmlParser("./pages/index.html")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	text, _ := lygo_io.ReadTextFromFile("./pages/index.html")
	parser, err = lygo_html.NewHtmlParser(text)
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	parser, err = lygo_html.NewHtmlParser("https://www.google.com")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	fmt.Println(parser.String())
}

func TestParser(t *testing.T) {

	parser, err := lygo_html.NewHtmlParser("./pages/index.html")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	path := parser.BaseUrl()
	fmt.Println("path", path)

	lang := parser.Lang()
	if len(lang) == 0 {
		t.Error("Expected lang")
		t.FailNow()
	}
	fmt.Println("lang", lang)

	title := parser.Title()
	if len(lang) == 0 {
		t.Error("Expected title")
		t.FailNow()
	}
	fmt.Println("title", title)

	fmt.Println("NODES:")
	parser.ForEach(func(node *html.Node) bool {
		fmt.Println("\t", node.Data, node.Type, node.Namespace, node.Attr)
		// fmt.Println(parser.InnerHtml(node))
		return false
	})

	fmt.Println("LINKS:")
	links := parser.GetLinkURLs()
	if len(links) == 0 {
		t.Error("Expected some links")
		t.FailNow()
	}
	for _, link := range links {
		fmt.Println("\t", link)
	}

	fmt.Println("PARAGRAPHS:")
	paragraphs := parser.Select("p")
	if len(paragraphs) == 0 {
		t.Error("Expected some paragraphs")
		t.FailNow()
	}
	for _, p := range paragraphs {
		fmt.Println("\t", parser.OuterHtml(p))
		fmt.Println("\t", parser.InnerHtml(p))
	}

	fmt.Println("KEYWORDS:")
	keywords := parser.MetaKeywords()
	if len(keywords) == 0 {
		t.Error("Expected some keywords")
		t.FailNow()
	}
	fmt.Println("\t", keywords)
}

func TestParserText(t *testing.T) {
	parser, err := lygo_html.NewHtmlParser("./pages/index.html")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	text := parser.TextAll()
	if len(text) == 0 {
		t.Error("Expected some text")
		t.FailNow()
	}
	fmt.Println("TEXT:")
	fmt.Println(text)
}

func TestParserSemantic(t *testing.T) {
	parser, err := lygo_html.NewHtmlParser("./pages/index.html")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	blocks := parser.SemanticBlocksAll()
	if len(blocks) == 0 {
		t.Error("Expected some blocks")
		t.FailNow()
	}
	fmt.Println("BLOCKS:", len(blocks))
	for _, block := range blocks {
		fmt.Println("\t", "title: ", block.Title)
		fmt.Println("\t", "lang: ", block.Lang)
		fmt.Println("\t", "block: ", block.Json())
	}
}

func TestParserPAth(t *testing.T) {

	paths:=[]string{"./pages/index.html", "https://gianangelogeminiani.me", "https://gianangelogeminiani.me/blog/", "https://gianangelogeminiani.me/intelligenza-artificiale-internet-of-things-e-blockchain-alla-portata-dei-tuoi-sistemi-informativi-ma-si-puo-fare/"}

	for _, path:=range paths{
		parser, err := lygo_html.NewHtmlParser(path)
		if nil != err {
			t.Error(err)
			t.FailNow()
		}
		fmt.Println("path", parser.Path())
		fmt.Println("root", parser.RootUrl())
		fmt.Println("base", parser.BaseUrl())
		fmt.Println("file", parser.FileName())
	}

}
