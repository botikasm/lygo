package test

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_json"
	"github.com/botikasm/lygo/ext/lygo_template"
	"testing"
)

func TestTemplate(t *testing.T) {

	model, err := lygo_json.ReadMapFromFile("./model.json")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	engine := lygo_template.NewTemplateEngine()
	err = engine.OpenModel("./tpl.txt")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("Paragraphs", len(engine.Paragraphs()))

	errs := engine.Render(model)
	if len(errs) > 0 {
		t.Error(errs)
		t.FailNow()
	}
	fmt.Println(engine)
	engine.SaveTo("./doc.txt")
}
