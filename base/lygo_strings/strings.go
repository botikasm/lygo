package lygo_strings

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_conv"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func TrimSpaces(slice []string) {
	Trim(slice, " ")
}

func Trim(slice []string, trimVal string) {
	for i := range slice {
		slice[i] = strings.Trim(slice[i], trimVal)
	}
}

func Concat(params ...interface{}) string {
	result := ""
	for _, v := range params {
		result += lygo_conv.ToString(v)
	}
	return result
}

func ConcatSep(separator string, params ...interface{}) string {
	result := ""
	strParams := lygo_conv.ToArrayOfString(params...)
	for _, value := range strParams {
		if len(result) > 0 {
			result += separator
		}
		result += value
	}
	return result
}

func ConcatTrimSep(separator string, params ...interface{}) string {
	result := ""
	for _, v := range params {
		value := strings.TrimSpace(lygo_conv.ToString(v))
		if len(value) > 0 {
			if len(result) > 0 {
				result += separator
			}
			result += value
		}
	}
	return result
}

func Format(s string, params ...interface{}) string {
	return fmt.Sprintf(strings.Replace(s, "%s", "%v", -1), params...)
}

func FormatValues(s string, params ...interface{}) string {
	return fmt.Sprintf(s, params...)
}

// Split using all rune in a string of separators
func Split(s string, seps string) []string {
	return strings.FieldsFunc(s, func(r rune) bool {
		for _, sep := range seps {
			if r == sep {
				return true
			}
		}
		return false
	})
}

func SplitTrim(s string, seps string, cutset string) []string {
	data := Split(s, seps)
	for i, item := range data {
		data[i] = strings.Trim(item, cutset)
	}
	return data
}

func SplitTrimSpace(s string, seps string) []string {
	return SplitTrim(s, seps, " ")
}

// get a substring
// @param s string The string
// @param start int Start index
// @param end int End index
func Sub(s string, start int, end int) string {
	if start < 0 || start > end {
		start = 0
	}
	if end > len(s) {
		end = len(s)
	}
	runes := []rune(s) // convert in rune to handle all characters.

	return string(runes[start:end])
}

//----------------------------------------------------------------------------------------------------------------------
//	C a m e l    C a s e
//----------------------------------------------------------------------------------------------------------------------

func CapitalizeAll(text string) string {
	return strings.Title(text)
}

func CapitalizeFirst(text string) string {
	if len(text) > 0 {
		words := Split(text, " ")
		if len(words) > 0 {
			words[0] = strings.Title(words[0])
			return ConcatSep(" ", words)
		}
	}
	return text
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
