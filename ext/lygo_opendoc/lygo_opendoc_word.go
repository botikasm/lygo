package lygo_opendoc

import (
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_crypto"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_reflect"
	"github.com/botikasm/lygo/base/lygo_rnd"
	"github.com/cbroglie/mustache"
	"github.com/unidoc/unioffice/document"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type OpenDocWord struct {

	//-- private --//
	filename string
	document *document.Document
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewOpenDocWord() *OpenDocWord {
	instance := new(OpenDocWord)

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------
func (instance *OpenDocWord) Filename() string {
	if nil != instance {
		return instance.filename
	}
	return ""
}

func (instance *OpenDocWord) Open(filename string) error {
	if nil != instance {
		doc, err := document.Open(filename)
		if nil != err {
			return err
		}
		instance.filename = filename
		instance.document = doc
	}
	return nil
}

// Open a file using as template and creates a temporary working file.
func (instance *OpenDocWord) OpenModel(filename string) error {
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

func (instance *OpenDocWord) Save() error {
	if nil != instance {
		if len(instance.filename) > 0 {
			lygo_paths.Mkdir(instance.filename)
			return instance.document.SaveToFile(instance.filename)
		}
	}
	return nil
}

func (instance *OpenDocWord) SaveTo(filename string) error {
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

func (instance *OpenDocWord) Render(context map[string]interface{}) []error {
	response := make([]error, 0)
	if nil != instance {

		// render paragraphs, except tables
		instance.scanParagraphs(instance.paragraphs(false), context)

		// render tables
		instance.scanTables(instance.document, context)
	}
	return response
}

func (instance *OpenDocWord) Paragraphs() []document.Paragraph {
	if nil != instance {
		return instance.paragraphs(true)
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *OpenDocWord) paragraphs(includeTable bool) []document.Paragraph {
	response := make([]document.Paragraph, 0)
	if nil != instance.document {
		doc := instance.document

		for _, h := range doc.Headers() {
			for _, p := range h.Paragraphs() {
				response = append(response, p)
			}
		}

		for _, p := range doc.Paragraphs() {
			response = append(response, p)
		}

		for _, sdt := range doc.StructuredDocumentTags() {
			for _, p := range sdt.Paragraphs() {
				response = append(response, p)
			}
		}

		if includeTable {
			for _, t := range doc.Tables() {
				for _, r := range t.Rows() {
					for _, c := range r.Cells() {
						for _, p := range c.Paragraphs() {
							response = append(response, p)
						}
					}
				}
			}
		}

		for _, f := range doc.Footers() {
			for _, p := range f.Paragraphs() {
				response = append(response, p)
			}
		}
	}
	return response
}

func (instance *OpenDocWord) scanParagraphs(paragraphs []document.Paragraph, context map[string]interface{}) []error {
	//paragraphs := instance.paragraphs(false)
	response := make([]error, 0)
	if len(paragraphs) > 0 {
		buffer := ""
		for _, p := range paragraphs {
			for _, r := range p.Runs() {
				text := r.Text()
				if len(text) == 0 {
					continue
				}
				if expStart(text) {
					// start buffered action
					buffer = text
					r.ClearContent()
				} else if expEnd(text) {
					// end buffered action
					buffer += text
					if !expIsTable(buffer) {
						err := instance.render(&r, buffer, context)
						if nil != err {
							response = append(response, err)
						}
					} else {
						// replace table
						r.ClearContent()
						r.AddText(buffer)
					}
					buffer = ""
				} else {
					if len(buffer) > 0 {
						// continue buffered action
						buffer += text
						r.ClearContent()
					} else {
						if expIsComplete(text) {
							if !expIsTable(text) {
								err := instance.render(&r, text, context)
								if nil != err {
									response = append(response, err)
								}
							} else {
								// replace table
								r.ClearContent()
								r.AddText(text)
							}
						}
					}
				}
			}
		}
	}
	return response
}

func (instance *OpenDocWord) scanTables(doc *document.Document, context map[string]interface{}) []error {
	response := make([]error, 0)
	nextRow := false
	for _, table := range doc.Tables() {
		for _, row := range table.Rows() {
			nextRow = false
			for _, c := range row.Cells() {
				if nextRow {
					break
				}
				for _, p := range c.Paragraphs() {
					if nextRow {
						break
					}
					for _, r := range p.Runs() {
						text := r.Text()
						if expIsComplete(text) {
							err := instance.renderTable(&table, row, context)
							if nil != err {
								response = append(response, err)
							}
							// jump to next rows
							nextRow = true
							break
						}
					}
				}
			}
		}
	}
	return response
}

func (instance *OpenDocWord) render(run *document.Run, text string, context map[string]interface{}) error {
	// STANDARD FIELD
	response, err := mustache.Render(text, context)
	if nil != err {
		return err
	}
	run.ClearContent()
	run.AddText(response)

	return nil
}

func (instance *OpenDocWord) renderTable(table *document.Table, row document.Row, context map[string]interface{}) error {
	cells := row.Cells()
	if len(cells) > 0 {
		// first cell must contain table data
		firstCell := cells[0]
		_, tableName, _ := parseTableName(textOfCell(firstCell, false))
		if len(tableName) > 0 {
			tableData := lygo_conv.ToArray(context[tableName])
			if nil != tableData {
				// get all expressions in cells and reset cells
				exprs := initCells(cells)
				// add rows to table
				idxFirstRow := len(table.Rows()) - 1
				addRows(table, len(tableData)-1, &row)
				// loop on all items in table
				for idxRow, item := range tableData {
					if nil != item {
						// loop on all cells
						for idxCell, expr := range exprs {
							if expIsComplete(expr) {
								_, _, fieldName := parseTableName(expr)
								if len(fieldName) > 0 {
									value := lygo_reflect.GetString(item, fieldName)
									var curRow document.Row
									if idxRow == 0 {
										// first row
										curRow = row
									} else {
										curRow = table.Rows()[idxFirstRow+idxRow]
									}
									cell := getCell(curRow, idxCell)
									setCellText(cell, value)
									// fmt.Println(idxRow, idxCell, len(curRow.Cells()), value)
								}
							}
						}
					}
				}
			}
		}
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	S T A T I C
//----------------------------------------------------------------------------------------------------------------------

func addRows(table *document.Table, count int, master *document.Row) {
	for i := 0; i < count; i++ {
		row := table.AddRow()
		addCells(row, len(master.Cells()))
		copyRowStyles(master, &row)
	}
}

func initCells(cells []document.Cell) []string {
	response := make([]string, 0)
	for _, c := range cells {
		response = append(response, textOfCell(c, true))
	}
	return response
}

func addCells(row document.Row, count int) {
	for i := 0; i < count; i++ {
		c := row.AddCell()
		p := c.AddParagraph()
		r := p.AddRun()
		r.ClearContent()
	}
}

func getCell(row document.Row, idx int) *document.Cell {
	if len(row.Cells()) > idx {
		return &row.Cells()[idx]
	}
	return nil
}

func setCellText(cell *document.Cell, text string) {
	var p document.Paragraph
	if len(cell.Paragraphs()) == 0 {
		p = cell.AddParagraph()
	} else {
		p = cell.Paragraphs()[0]
	}
	setParagraphText(p, text)
}

func copyCellStyles(from, to document.Cell) {
	for _, p1 := range from.Paragraphs() {
		style := p1.Style()
		if len(style) > 0 {
			for _, p2 := range to.Paragraphs() {
				p2.SetStyle(style)
			}
		}
		for _, r := range p1.Runs() {
			runProps := r.Properties()
			for _, p2 := range to.Paragraphs() {
				for _, r2 := range p2.Runs() {
					r2.Properties().SetBold(runProps.IsBold())
				}
			}
		}
	}
}

func copyRowStyles(from, to *document.Row) {
	c1 := from.Cells()
	c2 := to.Cells()
	if len(c1) == len(c2) {
		for i, cc1 := range c1 {
			copyCellStyles(cc1, c2[i])
		}
	}
}

func setParagraphText(p document.Paragraph, text string) {
	for i, r := range p.Runs() {
		r.ClearContent()
		if len(text) > 0 && i == 0 {
			r.AddText(text)
		}
	}
}

func textOfCell(c document.Cell, clear bool) string {
	response := ""
	for _, p := range c.Paragraphs() {
		for _, r := range p.Runs() {
			response += r.Text()
			if clear {
				r.ClearContent()
			}
		}
	}
	return response
}

func parseTableName(text string) (string, string, string) {
	text = strings.ReplaceAll(text, "[", "|")
	text = strings.ReplaceAll(text, "]", "") // {{table|name}}
	tableName := text[strings.Index(text, "{{")+2 : strings.Index(text, "|")]
	fieldName := text[strings.Index(text, "|")+1 : strings.Index(text, "}}")]
	return text, tableName, fieldName
}

func expIsTable(text string) bool {
	return strings.Index(text, "]}}") > -1
}

func expIsComplete(text string) bool {
	return strings.Index(text, "{{") > -1 && strings.Index(text, "}}") > -1
}

func expStart(text string) bool {
	return !expIsComplete(text) && strings.Index(text, "{{") > -1
}

func expEnd(text string) bool {
	return !expIsComplete(text) && strings.Index(text, "}}") > -1
}
