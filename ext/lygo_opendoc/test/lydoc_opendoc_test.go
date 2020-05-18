package lygo_opendoc

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_exec"
	"github.com/botikasm/lygo/ext/lygo_opendoc"
	"github.com/cbroglie/mustache"
	"github.com/unidoc/unioffice/color"
	"github.com/unidoc/unioffice/document"
	"github.com/unidoc/unioffice/measurement"
	"github.com/unidoc/unioffice/schema/soo/wml"
	"testing"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func TestOpenDocTemplate(t *testing.T) {
	doc := lygo_opendoc.NewOpenDocWord()
	err := doc.OpenModel("./tpl.docx")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}
	fmt.Println("FILE OPENED AND WORKING ON:", doc.Filename())

	context := map[string]interface{}{
		"header": map[string]interface{}{
			"company": "Gemini",
			"phone":   "+393454643543",
		},
		"name":  "John",
		"email": "angelo.geminiani@gmail.com",
		"table": []map[string]interface{}{
			{
				"name":    "Carlo",
				"surname": "Martello",
			},
			{
				"name":    "Mario",
				"surname": "Cartello",
			},
		},
	}
	errs := doc.Render(context)
	if len(errs) > 0 {
		fmt.Println(len(errs), errs)
		t.Error(errs[0])
		t.FailNow()
	}
	err = doc.SaveTo("./tpl-doc1.docx")
	if nil != err {
		t.Error(err)
		t.FailNow()
	}

	lygo_exec.Open("./tpl-doc1.docx")
}

func TestSimpleWord(t *testing.T) {

	doc := document.New()

	para := doc.AddParagraph()
	run := para.AddRun()
	para.SetStyle("Title")
	run.AddText("Simple Document Formatting")

	para = doc.AddParagraph()
	para.SetStyle("Heading1")
	run = para.AddRun()
	run.AddText("Some Heading Text")

	para = doc.AddParagraph()
	para.SetStyle("Heading2")
	run = para.AddRun()
	run.AddText("Some Heading Text")

	para = doc.AddParagraph()
	para.SetStyle("Heading3")
	run = para.AddRun()
	run.AddText("Some Heading Text")

	para = doc.AddParagraph()
	para.Properties().SetFirstLineIndent(0.5 * measurement.Inch)

	run = para.AddRun()
	run.AddText("A run is a string of characters with the same formatting. ")

	run = para.AddRun()
	run.Properties().SetBold(true)
	run.Properties().SetFontFamily("Courier")
	run.Properties().SetSize(15)
	run.Properties().SetColor(color.Red)
	run.AddText("Multiple runs with different formatting can exist in the same paragraph. ")

	run = para.AddRun()
	run.AddText("Adding breaks to a run will insert line breaks after the run. ")
	run.AddBreak()
	run.AddBreak()

	run = createParaRun(doc, "Runs support styling options:")

	run = createParaRun(doc, "small caps")
	run.Properties().SetSmallCaps(true)

	run = createParaRun(doc, "strike")
	run.Properties().SetStrikeThrough(true)

	run = createParaRun(doc, "double strike")
	run.Properties().SetDoubleStrikeThrough(true)

	run = createParaRun(doc, "outline")
	run.Properties().SetOutline(true)

	run = createParaRun(doc, "emboss")
	run.Properties().SetEmboss(true)

	run = createParaRun(doc, "shadow")
	run.Properties().SetShadow(true)

	run = createParaRun(doc, "imprint")
	run.Properties().SetImprint(true)

	run = createParaRun(doc, "highlighting")
	run.Properties().SetHighlight(wml.ST_HighlightColorYellow)

	run = createParaRun(doc, "underline")
	run.Properties().SetUnderline(wml.ST_UnderlineWavyDouble, color.Red)

	run = createParaRun(doc, "text effects")
	run.Properties().SetEffect(wml.ST_TextEffectAntsRed)

	nd := doc.Numbering.Definitions()[0]

	for i := 1; i < 5; i++ {
		p := doc.AddParagraph()
		p.SetNumberingLevel(i - 1)
		p.SetNumberingDefinition(nd)
		run := p.AddRun()
		run.AddText(fmt.Sprintf("Level %d", i))
	}
	doc.SaveToFile("simple.docx")

}

func TestTemplateEngine(t *testing.T) {
	tpl := "{{#users}}\n{{name}} {{surname}}\n{{/users}}"
	model := map[string]interface{}{
		"users":[]map[string]interface{}{
			{
				"name":"Angelo",
				"surname":"Geminiani",
			},
			{
				"name":"Luca",
				"surname":"",
			},
			{
				"name":"Andrea",
				"surname":"",
			},
			{
				"name":"Alessandro",
				"surname":"",
			},
		},
	}
	out, err := mustache.Render(tpl, model)
	if nil!=err{
		t.Error(err)
		t.FailNow()
	}
	fmt.Println(out)
}


//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func createParaRun(doc *document.Document, s string) document.Run {
	para := doc.AddParagraph()
	run := para.AddRun()
	run.AddText(s)
	return run
}
