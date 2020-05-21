package lygo_template

import (
	"bytes"
	"github.com/botikasm/lygo/base/lygo_array"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_crypto"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_reflect"
	"github.com/botikasm/lygo/base/lygo_rnd"
	"github.com/botikasm/lygo/base/lygo_strings"
	"github.com/cbroglie/mustache"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type TemplateEngine struct {

	//-- private --//
	filename string
	source   string
	render   string
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewTemplateEngine() *TemplateEngine {
	instance := new(TemplateEngine)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *TemplateEngine) String() string {
	if nil != instance {
		return instance.render
	}
	return ""
}

func (instance *TemplateEngine) Filename() string {
	if nil != instance {
		return instance.filename
	}
	return ""
}

func (instance *TemplateEngine) Open(filename string) error {
	if nil != instance {
		text, err := lygo_io.ReadTextFromFile(filename)
		if nil != err {
			return err
		}
		instance.filename = filename
		instance.source = text
		instance.render = text
	}
	return nil
}

// Open a file using as template and creates a temporary working file.
func (instance *TemplateEngine) OpenModel(filename string) error {
	if nil != instance {
		err := instance.Open(filename)
		if nil != err {
			return err
		}
		// save temp file
		name := lygo_paths.FileName(instance.filename, false)
		uuid := lygo_crypto.MD5(lygo_rnd.Uuid())
		instance.filename = lygo_paths.Concat(lygo_paths.GetTempRoot(), name+"."+uuid+lygo_paths.Extension(instance.filename))
		return instance.Save()
	}
	return nil
}

func (instance *TemplateEngine) Save() error {
	if nil != instance {
		if len(instance.filename) > 0 {
			lygo_paths.Mkdir(instance.filename)
			return instance.saveToFile(instance.filename)
		}
	}
	return nil
}

func (instance *TemplateEngine) SaveTo(filename string) error {
	if nil != instance {
		if len(filename) > 0 {
			if lygo_paths.IsTemp(instance.filename) {
				lygo_io.Remove(instance.filename)
			}
			instance.filename = filename
			return instance.Save()
		}
	}
	return nil
}

func (instance *TemplateEngine) Render(context map[string]interface{}) []error {
	response := make([]error, 0)
	if nil != instance {
		text, errs := instance.renderAll(instance.source, context)
		if len(errs) > 0 {
			response = append(response, errs...)
		} else {
			instance.render = text
		}
	}
	return response
}

func (instance *TemplateEngine) Paragraphs() []string {
	if nil != instance {
		return instance.paragraphs(instance.render)
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *TemplateEngine) saveToFile(filename string) error {
	if nil != instance {
		_, err := lygo_io.WriteTextToFile(instance.render, filename)
		return err
	}
	return nil
}

func (instance *TemplateEngine) paragraphs(text string) []string {
	text = strings.ReplaceAll(text, "\r", "")
	return strings.Split(text, "\n")
}

func (instance *TemplateEngine) renderAll(text string, context map[string]interface{}) (string, []error) {
	errResponse := make([]error, 0)
	var buffer bytes.Buffer
	paragraphs := instance.paragraphs(text)
	for _, paragraph := range paragraphs {
		if expIsComplete(paragraph) {
			text, err := instance.renderParagraph(paragraph, context)
			if nil != err {
				buffer.WriteString(paragraph + "\n")
				errResponse = append(errResponse, err)
			} else {
				buffer.WriteString(text + "\n")
			}
		} else {
			buffer.WriteString(paragraph + "\n")
		}
	}

	return buffer.String(), errResponse
}

func (instance *TemplateEngine) renderParagraph(paragraph string, context map[string]interface{}) (string, error) {
	if expIsTable(paragraph) {
		// render a table
		return instance.renderTable(paragraph, context)
	}
	return mustache.Render(paragraph, context)
}

func (instance *TemplateEngine) renderTable(paragraph string, context map[string]interface{}) (string, error) {
	var buffer bytes.Buffer
	exps := parseTableFields(paragraph)
	if len(exps) > 0 {
		// adjust paragraph replacing table fields
		for _, exp := range exps {
			_, _, f, _ := parseTableName(exp)
			paragraph = strings.ReplaceAll(paragraph, exp, f)
		}
		_, tableName, fieldName, _ := parseTableName(exps[0])
		if len(tableName) > 0 && len(fieldName) > 0 {
			tableData := lygo_conv.ToArray(context[tableName])
			if nil != tableData {
				for i, item := range tableData {
					// build a model for current row
					model := map[string]interface{}{}
					for _, exp := range exps {
						_, _, f, cl := parseTableName(exp)
						value := lygo_reflect.Get(item, f)
						if v, b := value.(string); b && cl > 0 {
							value = lygo_strings.FillRight(v, cl, ' ')
						}
						model[f] = value
					}
					// render row
					text, err := mustache.Render(paragraph, model)
					if nil != err {
						return "", err
					}
					if len(text) > 0 {
						buffer.WriteString(text)
						if i < len(tableData)-1 {
							buffer.WriteString("\n")
						}
					}
				}
			}
		}
	}
	return buffer.String(), nil
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func parseTableName(text string) (string, string, string, int) {
	text = strings.ReplaceAll(text, "[", "|")
	text = strings.ReplaceAll(text, "]", ":") // {{table|name}}
	tokens := strings.Split(text, "|")
	tableName := lygo_array.GetAt(tokens, 0, "").(string)
	fieldNames := strings.Split(lygo_array.GetAt(tokens, 1, "").(string), ":")
	fieldName := fieldNames[0]
	columnLen := lygo_conv.ToIntDef(lygo_array.GetAt(fieldNames, 1, "0"), 0)

	return text, tableName, fieldName, columnLen
}

func parseTableFields(text string) []string {
	response := make([]string, 0)
	text = strings.ReplaceAll(text, " ", "")
	text = strings.ReplaceAll(text, "\t", "")
	text = strings.ReplaceAll(text, "}}{{", "|")
	text = strings.ReplaceAll(text, "}}", "")
	text = strings.ReplaceAll(text, "{{", "")

	exps := strings.Split(text, "|")
	for _, exp := range exps {
		if len(exp) > 0 {
			response = append(response, exp)
		}
	}
	return response
}

func expIsTable(text string) bool {
	return expIsComplete(text) && strings.Index(text, "[") > -1 && strings.Index(text, "]") > -1
}

func expIsComplete(text string) bool {
	return strings.Index(text, "{{") > -1 && strings.Index(text, "}}") > -1
}
