package lygo_strings

import (
	"bytes"
	"fmt"
	"github.com/botikasm/lygo/base/lygo_conv"
	"regexp"
	"strconv"
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

func Clear(text string) string {
	var buf bytes.Buffer
	lines := strings.Split(text, "\n")
	count := 0
	for _, line:=range lines{
		space := regexp.MustCompile(`\s+`)
		s := strings.TrimSpace(space.ReplaceAllString(line, " "))
		if len(s)>0{
			if count>0{
				buf.WriteString("\n")
			}
			buf.WriteString(strings.TrimSpace(s))
			count++
		}
	}
	return buf.String()
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

func SplitAndGetAt(s string, seps string, index int) string {
	tokens := Split(s, seps)
	if len(tokens) > index {
		return tokens[index]
	}
	return ""
}

// get a substring
// @param s string The string
// @param start int Start index
// @param end int End index
func Sub(s string, start int, end int) string {
	runes := []rune(s) // convert in rune to handle all characters.
	if start < 0 || start > end {
		start = 0
	}
	if end > len(runes) {
		end = len(runes)
	}


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
//	p a d d i n g
//----------------------------------------------------------------------------------------------------------------------

func FillLeft(text string, l int, r rune) string {
	if len(text) == l {
		return text
	} else if len(text) < l {
		return fmt.Sprintf("%"+string(r)+strconv.Itoa(l)+"s", text)
	}
	return text[:l]
}

func FillRight(text string, l int, r rune) string {
	if len(text) == l {
		return text
	} else if len(text) < l {
		return text + strings.Repeat(string(r), l-len(text))
	}
	return text[:l]
}

func FillLeftBytes(bytes []byte, l int, r rune) []byte {
	return []byte(FillLeft(string(bytes), l, r))
}

func FillLeftZero(text string, l int) string {
	return FillLeft(text, l, '0')
}

func FillLeftBytesZero(bytes []byte, l int) []byte {
	return []byte(FillLeftZero(string(bytes), l))
}

func FillRightZero(text string, l int) string {
	return FillRight(text, l, '0')
}

func FillRightBytes(bytes []byte, l int, r rune) []byte {
	return []byte(FillRight(string(bytes), l, r))
}

func FillRightBytesZero(bytes []byte, l int) []byte {
	return []byte(FillRight(string(bytes), l, '0'))
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
