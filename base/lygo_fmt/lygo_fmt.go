package lygo_fmt

import (
	"strings"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	D A T E
//----------------------------------------------------------------------------------------------------------------------

func FormatDate(dt time.Time, pattern string) string {
	return dt.Format(toGoLayout(pattern))
}

func ParseDate(dt string, pattern string) (time.Time, error) {
	return time.Parse(toGoLayout(pattern), dt)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func toGoLayout(pattern string) string {
	response := strings.ReplaceAll(pattern, "yyyy", "2006")
	response = strings.ReplaceAll(response, "MM", "01")
	response = strings.ReplaceAll(response, "dd", "02")
	response = strings.ReplaceAll(response, "HH", "15")
	response = strings.ReplaceAll(response, "mm", "04")
	response = strings.ReplaceAll(response, "ss", "05")
	response = strings.ReplaceAll(response, "Z", "-0700")
	return response
}
