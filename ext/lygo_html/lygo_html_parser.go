package lygo_html

import (
	"bytes"
	"github.com/botikasm/lygo/base/lygo_array"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_strings"
	"github.com/botikasm/lygo/ext/lygo_http/lygo_http_client"
	"github.com/botikasm/lygo/ext/lygo_nlp/lygo_nlpdetect"
	"golang.org/x/net/html"
	"golang.org/x/text/language"
	"io"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	SemanticBlock
//----------------------------------------------------------------------------------------------------------------------

type SemanticBlock struct {
	Lang     string // detected language (may be different from page lang)
	Level    int    // title level. 0 is when a block is free text with no title
	Title    string
	Body     bytes.Buffer
	Keywords []string
}

func (instance *SemanticBlock) String() string {
	if nil != instance {
		var buf bytes.Buffer
		if len(instance.Title) > 0 {
			buf.WriteString(instance.Title)
			buf.WriteString("\n")
		}
		if instance.Body.Len() > 0 {
			buf.Write(instance.Body.Bytes())
		}
		return buf.String()
	}
	return ""
}

//----------------------------------------------------------------------------------------------------------------------
//	HtmlParser
//----------------------------------------------------------------------------------------------------------------------

type HtmlParser struct {
	html *html.Node
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewHtmlParser(input interface{}) (*HtmlParser, error) {
	var doc *html.Node
	var err error

	// parse input data
	if v, b := input.(string); b {
		doc, err = parseString(v)
	} else if v, b := input.(io.Reader); b {
		doc, err = parse(v)
	}

	if nil != err {
		return nil, err
	} else {
		instance := new(HtmlParser)
		instance.html = doc
		return instance, nil
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *HtmlParser) String() string {
	if nil != instance {
		return renderNode(instance.html)
	}
	return ""
}

func (instance *HtmlParser) Lang() string {
	lang := ""
	if nil != instance && nil != instance.html {
		forEach(instance.html, func(node *html.Node) bool {
			lang = getAttr(node, "lang")
			return len(lang) > 0 // next node?
		})
	}
	return lang
}

func (instance *HtmlParser) Title() string {
	title := ""
	if nil != instance && nil != instance.html {
		forEach(instance.html, func(node *html.Node) bool {
			if strings.ToLower(node.Data) == "title" {
				title = instance.InnerHtml(node)
			}
			return len(title) > 0 // next node?
		})
	}
	return title
}

func (instance *HtmlParser) MetaTitle() string {
	value := ""
	if nil != instance && nil != instance.html {
		value = instance.GetMetaContent("title")
	}
	return value
}

func (instance *HtmlParser) MetaDescription() string {
	value := ""
	if nil != instance && nil != instance.html {
		value = instance.GetMetaContent("description")
	}
	return value
}

func (instance *HtmlParser) MetaAuthor() string {
	value := ""
	if nil != instance && nil != instance.html {
		value = instance.GetMetaContent("author")
	}
	return value
}

func (instance *HtmlParser) MetaKeywords() []string {
	if nil != instance && nil != instance.html {
		return lygo_strings.SplitTrimSpace(instance.GetMetaContent("keywords"), ",")
	}
	return []string{}
}

func (instance *HtmlParser) GetText() []string {
	if nil != instance && nil != instance.html {
		return lygo_strings.SplitTrimSpace(instance.GetMetaContent("keywords"), ",")
	}
	return []string{}
}

func (instance *HtmlParser) GetMetaContent(name string) string {
	value := ""
	if nil != instance && nil != instance.html {
		forEach(instance.html, func(node *html.Node) bool {
			if strings.ToLower(node.Data) == "meta" {
				attrName := getAttr(node, "name")
				if attrName == name {
					value = getAttr(node, "content")
				}
			}
			return len(value) > 0 // next node?
		})
	}
	return value
}

func (instance *HtmlParser) TextAll() string {
	if nil != instance {
		return text(instance.html)
	}
	return ""
}

func (instance *HtmlParser) Text(node *html.Node) string {
	if nil != instance {
		return text(node)
	}
	return ""
}

func (instance *HtmlParser) SemanticBlocksAll() []*SemanticBlock {
	if nil != instance {
		return semantic(instance.Lang(), instance.html)
	}
	return make([]*SemanticBlock, 0)
}

func (instance *HtmlParser) SemanticBlocks(node *html.Node) []*SemanticBlock {
	if nil != instance {
		return semantic(instance.Lang(), node)
	}
	return make([]*SemanticBlock, 0)
}

func (instance *HtmlParser) OuterHtml(n *html.Node) string {
	if nil != instance {
		return renderNode(n)
	}
	return ""
}

func (instance *HtmlParser) InnerHtml(n *html.Node) string {
	if nil != instance {
		return renderNode(n.FirstChild)
	}
	return ""
}

func (instance *HtmlParser) ForEach(callback func(node *html.Node) bool) {
	if nil != instance && nil != instance.html && nil != callback {
		forEach(instance.html, callback)
	}
}

func (instance *HtmlParser) GelLinks() []*html.Node {
	if nil != instance && nil != instance.html {
		return queryNodes(instance.html, "a")
	}
	return []*html.Node{}
}

func (instance *HtmlParser) GelLinkURLs() []string {
	response := make([]string, 0)
	if nil != instance && nil != instance.html {
		links := queryNodes(instance.html, "a")
		for _, link := range links {
			href := getAttr(link, "href")
			if len(href) > 0 {
				response = append(response, href)
			}
		}
	}
	return response
}

func (instance *HtmlParser) GeNodeAttributes(nodes []*html.Node) []map[string]string {
	response := make([]map[string]string, 0)
	if nil != instance && nil != instance.html {
		for _, node := range nodes {
			m := map[string]string{
				"tag": node.Data,
			}
			for _, attr := range node.Attr {
				m[attr.Key] = attr.Val
			}
			response = append(response, m)
		}
	}
	return response
}

func (instance *HtmlParser) Select(selector string) []*html.Node {
	if nil != instance && nil != instance.html {
		return queryNodes(instance.html, selector)
	}
	return []*html.Node{}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func parse(input io.Reader) (*html.Node, error) {
	return html.Parse(input)
}

func parseString(input string) (*html.Node, error) {
	if strings.Index(input, "http") == 0 {
		return parseURL(input)
	} else if b, err := lygo_paths.IsFile(input); b && nil == err {
		return parseFile(input)
	}
	return html.Parse(strings.NewReader(input))
}

func parseBytes(input []byte) (*html.Node, error) {
	return html.Parse(bytes.NewReader(input))
}

func parseFile(input string) (*html.Node, error) {
	text, err := lygo_io.ReadTextFromFile(input)
	if nil != err {
		return nil, err
	}
	return html.Parse(strings.NewReader(text))
}

func parseURL(input string) (*html.Node, error) {
	client := lygo_http_client.NewHttpClient()
	_, data, err := client.Get(input)
	if nil != err {
		return nil, err
	}
	return parseBytes(data)
}

func renderNode(n *html.Node) string {
	if nil != n {
		var buf bytes.Buffer
		w := io.Writer(&buf)
		_ = html.Render(w, n)
		return buf.String()
	}
	return ""
}

func isTitle(node *html.Node) (bool, int) {
	tag := strings.ToLower(node.Data)
	if len(tag) == 2 && strings.Index(tag, "h") == 0 {
		level := tag[1]
		return true, lygo_conv.ToInt(string(level))
	}
	return false, -1
}

func text(node *html.Node) string {
	var buf bytes.Buffer
	forEach(node, func(node *html.Node) bool {
		if node.Data != "title" && nil != node.FirstChild && node.FirstChild.Type == html.TextNode {
			text := lygo_strings.Clear(renderNode(node.FirstChild))
			if len(text) > 0 {
				buf.WriteString(text + "\n")
			}
		}
		return false // next
	})
	return buf.String()
}

func unescape(text string) string {
	return html.UnescapeString(lygo_strings.Clear(text))
}

func detectLanguage(text string) string {
	txt := strings.ReplaceAll(text, "\n", " ")
	i18n := lygo_nlpdetect.DetectOne(strings.ToLower(txt))
	t, err := language.Parse(i18n.Code)
	if nil == err {
		return t.String()
	}
	return ""
}

func semantic(lang string, node *html.Node) []*SemanticBlock {
	response := make([]*SemanticBlock, 0)
	tmpBlock := new(SemanticBlock)
	tmpBlock.Lang = lang
	forEach(node, func(node *html.Node) bool {
		if node.Data != "title" && nil != node.FirstChild && node.FirstChild.Type == html.TextNode {
			if b, level := isTitle(node); b {
				// new block
				if tmpBlock.Body.Len() > 0 || len(tmpBlock.Title) > 0 {
					response = append(response, tmpBlock)
				}
				tmpBlock = new(SemanticBlock)
				tmpBlock.Lang = lang
				tmpBlock.Title = unescape(renderNode(node.FirstChild))
				tmpBlock.Level = level
			} else {
				text := unescape(renderNode(node.FirstChild))
				if len(text) > 0 {
					tmpBlock.Body.WriteString(text + "\n")
				}
			}
		}
		return false // next
	})

	if tmpBlock.Body.Len() > 0 {
		response = append(response, tmpBlock)
	}

	// semantic check
	for _, block := range response {
		// detect language
		lang := detectLanguage(block.String())
		if len(lang) > 0 {
			block.Lang = lang
		}
		// detect keywords

	}

	return response
}

func forEach(node *html.Node, callback func(node *html.Node) bool) {
	if nil != node {
		for child := node.FirstChild; child != nil; child = child.NextSibling {
			if child.Type == html.ElementNode || nil != child.FirstChild {
				exit := callback(child)
				if exit {
					break
				}
				forEach(child, callback)
			}
		}
	}
}

func getAttr(node *html.Node, name string) string {
	if nil != node && len(node.Attr) > 0 {
		for _, attr := range node.Attr {
			if strings.ToLower(attr.Key) == strings.ToLower(name) {
				return attr.Val
			}
		}
	}
	return ""
}

func queryNodes(root *html.Node, selector string) []*html.Node {
	response := make([]*html.Node, 0)
	forEach(root, func(node *html.Node) bool {
		if matches(node, selector) {
			response = append(response, node)
		}
		return false // next node
	})
	return response
}

func matches(node *html.Node, selector string) bool {
	if strings.Index(selector, ".") == 0 {
		// class matching
		className := strings.Replace(selector, ".", "", 1)
		classes := strings.Split(getAttr(node, "class"), " ")
		return lygo_array.IndexOf(className, classes) > -1
	} else {
		// tag name matching
		return node.Data == selector
	}
}
