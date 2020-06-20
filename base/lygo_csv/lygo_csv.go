package lygo_csv

import (
	"encoding/csv"
	"fmt"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------


//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type CsvOptions struct {
	Comma          string `json:"comma"`
	Comment        string `json:"comment"`
	FirstRowHeader bool   `json:"first_row_header"`
}

func NewCsvOptions(comma string, comment string, firstRowHeader bool) *CsvOptions {
	return &CsvOptions{
		Comma:          comma,
		Comment:        comment,
		FirstRowHeader: firstRowHeader,
	}
}

func NewCsvOptionsDefaults() *CsvOptions {
	return &CsvOptions{
		Comma:          ";",
		Comment:        "#",
		FirstRowHeader: true,
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func ReadAll(in string, options *CsvOptions) (response []map[string]string, err error) {

	r := csv.NewReader(strings.NewReader(in))
	setOptions(r, options)

	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}

	response = make([]map[string]string, 0)
	headers := buildHeaders(&records, options)
	for _, row := range records {
		item := make(map[string]string)
		for i, value := range row {
			if len(headers) > i {
				item[headers[i]] = value
			}
		}
		response = append(response, item)
	}

	return response, err
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func setOptions(r *csv.Reader, options *CsvOptions) {
	if nil != r && nil != options {
		if len(options.Comma) == 1 {
			r.Comma = []rune(options.Comma)[0]
		}
		if len(options.Comment) == 1 {
			r.Comment = []rune(options.Comment)[0]
		}
	}
}

func buildHeaders(records *[][]string, options *CsvOptions) []string {
	headers := make([]string, 0)
	if options.FirstRowHeader && len(*records) > 1 {
		headers = (*records)[0]
		*records = (*records)[1:][:]
	} else {
		for i := 0; i < len(*records); i++ {
			headers = append(headers, fmt.Sprintf("field_%v", i))
		}
	}

	return headers
}
