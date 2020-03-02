package lygo_regex

import (
	"encoding/json"
	"fmt"
	"github.com/botikasm/lygo/base/lygo_strings"
	"regexp"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

// Date finds all date strings
func Date(text string) []string {
	return matchString(text, DateRegex)
}

// Time finds all time strings
func Time(text string) []string {
	return matchString(text, TimeRegex)
}

// Phones finds all phone numbers
func Phones(text string) []string {
	return matchString(text, PhoneRegex)
}

// PhonesWithExts finds all phone numbers with ext
func PhonesWithExts(text string) []string {
	return matchString(text, PhonesWithExtsRegex)
}

// Links finds all link strings
func Links(text string) []string {
	return matchString(text, LinkRegex)
}

// Emails finds all email strings
func Emails(text string) []string {
	return matchString(text, EmailRegex)
}

// IPv4s finds all IPv4 addresses
func IPv4s(text string) []string {
	return matchString(text, IPv4Regex)
}

// IPv6s finds all IPv6 addresses
func IPv6s(text string) []string {
	return matchString(text, IPv6Regex)
}

// IPs finds all IP addresses (both IPv4 and IPv6)
func IPs(text string) []string {
	return matchString(text, IPRegex)
}

// NotKnownPorts finds all not-known port numbers
func NotKnownPorts(text string) []string {
	return matchString(text, NotKnownPortRegex)
}

// Prices finds all price strings
func Prices(text string) []string {
	array := matchString(text, PriceRegex)
	lygo_strings.TrimSpaces(array)
	return array
}

// HexColors finds all hex color values
func HexColors(text string) []string {
	return matchString(text, HexColorRegex)
}

// CreditCards finds all credit card numbers
func CreditCards(text string) []string {
	return matchString(text, CreditCardRegex)
}

// BtcAddresses finds all bitcoin addresses
func BtcAddresses(text string) []string {
	return matchString(text, BtcAddressRegex)
}

// StreetAddresses finds all street addresses
func StreetAddresses(text string) []string {
	return matchString(text, StreetAddressRegex)
}

// ZipCodes finds all zip codes
func ZipCodes(text string) []string {
	return matchString(text, ZipCodeRegex)
}

// PoBoxes finds all po-box strings
func PoBoxes(text string) []string {
	return matchString(text, PoBoxRegex)
}

// SSNs finds all SSN strings
func SSNs(text string) []string {
	return matchString(text, SSNRegex)
}

// MD5Hexes finds all MD5 hex strings
func MD5Hexes(text string) []string {
	return matchString(text, MD5HexRegex)
}

// SHA1Hexes finds all SHA1 hex strings
func SHA1Hexes(text string) []string {
	return matchString(text, SHA1HexRegex)
}

// SHA256Hexes finds all SHA256 hex strings
func SHA256Hexes(text string) []string {
	return matchString(text, SHA256HexRegex)
}

// GUIDs finds all GUID strings
func GUIDs(text string) []string {
	return matchString(text, GUIDRegex)
}

// ISBN13s finds all ISBN13 strings
func ISBN13s(text string) []string {
	return matchString(text, ISBN13Regex)
}

// ISBN10s finds all ISBN10 strings
func ISBN10s(text string) []string {
	return matchString(text, ISBN10Regex)
}

// VISACreditCards finds all VISA credit card numbers
func VISACreditCards(text string) []string {
	return matchString(text, VISACreditCardRegex)
}

// MCCreditCards finds all MasterCard credit card numbers
func MCCreditCards(text string) []string {
	return matchString(text, MCCreditCardRegex)
}

// MACAddresses finds all MAC addresses
func MACAddresses(text string) []string {
	return matchString(text, MACAddressRegex)
}

// IBANs finds all IBAN strings
func IBANs(text string) []string {
	return matchString(text, IBANRegex)
}

// GitRepos finds all git repository addresses which have protocol prefix
func GitRepos(text string) []string {
	return matchString(text, GitRepoRegex)
}

func Numbers(text string) []string {
	return matchString(text, NumbersRegex)
}

//----------------------------------------------------------------------------------------------------------------------
//	v a l i d a t i o n
//----------------------------------------------------------------------------------------------------------------------

func IsValidEmail(text string) bool {
	return len(Emails(text)) == 1
}

func IsValidJsonObject(text string) bool {
	var js map[string]interface{}
	return json.Unmarshal([]byte(text), &js) == nil
}

func IsValidJsonArray(text string) bool {
	var js []map[string]interface{}
	return json.Unmarshal([]byte(text), &js) == nil
}

//----------------------------------------------------------------------------------------------------------------------
//	w i l d c a r d    l o o k u p
//----------------------------------------------------------------------------------------------------------------------

func WildcardMatchAll(text, expression string) ([]string, [][]int) {
	exp := toRegexp(expression)
	return matchAll(text, exp)
}

func WildcardMatch(text, expression string) []string {
	exp := toRegexp(expression)
	return matchString(text, exp)
}

func WildcardMatchIndex(text, expression string) [][]int {
	exp := toRegexp(expression)
	return matchIndex(text, exp)
}

func WildcardMatchBetween(text string, offset int, patternStart string, patternEnd string, cutset string) []string {

	expStart := toRegexp(patternStart)
	expEnd := toRegexp(patternEnd)

	return matchBetween(text, offset, expStart, expEnd, cutset)
}

// Return index array of matching expression in a text starting search from offset position
// @param text string. "hello humanity!!"
// @param pattern string "hu?an*"
// @param offset int number of characters to exclude from search
// @return []int
func WildcardIndex(text string, pattern string, offset int) []int {
	regex := toRegexp(pattern)
	return index(text, regex, offset)
}

// Return array of pair index:word_len  of matching expression in a text
// @param text string. "hello humanity!!"
// @param pattern string "hu?an*"
// @return [][]int ex: [[12,3], [22,4]]
func WildcardIndexLenPair(text string, pattern string, offset int) [][]int {
	regex := toRegexp(pattern)
	return indexLenPair(text, regex, offset)
}

//----------------------------------------------------------------------------------------------------------------------
//	e x p    l o o k u p
//----------------------------------------------------------------------------------------------------------------------

func MatchAll(text, expression string) ([]string, [][]int) {
	exp := regexp.MustCompile(expression)
	return matchAll(text, exp)
}

func Match(text, expression string) []string {
	exp := regexp.MustCompile(expression)
	return matchString(text, exp)
}

func MatchIndex(text, expression string) [][]int {
	exp := regexp.MustCompile(expression)
	return matchIndex(text, exp)
}

func MatchBetween(text string, offset int, patternStart string, patternEnd string, cutset string) []string {
	expStart := regexp.MustCompile(patternStart)
	expEnd := regexp.MustCompile(patternEnd)

	return matchBetween(text, offset, expStart, expEnd, cutset)
}

func Index(text string, pattern string, offset int) []int {
	regex := regexp.MustCompile(pattern)
	return index(text, regex, offset)
}

func IndexLenPair(text string, pattern string, offset int) [][]int {
	regex := regexp.MustCompile(pattern)
	return indexLenPair(text, regex, offset)
}

//----------------------------------------------------------------------------------------------------------------------
//	s p l i t
//----------------------------------------------------------------------------------------------------------------------

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func toRegexp(wildCardExpr string) *regexp.Regexp {
	if wildCardExpr == "\n" {
		return regexp.MustCompile(wildCardExpr)
	} else {
		prefix := "(\\b"
		suffix := "\\b)"

		// replace spaces
		wildCardExpr = strings.Replace(wildCardExpr, " ", "\\W", -1)

		// escape dot
		wildCardExpr = strings.Replace(wildCardExpr, ".", "\\.", -1)

		// replace ?
		wildCardExpr = strings.Replace(wildCardExpr, "?", "(?:(?:.|\n))?", -1)

		// replace *
		wildCardExpr = strings.Replace(wildCardExpr, "*", ".*?", -1)

		r := fmt.Sprintf("%s%s%s", prefix, wildCardExpr, suffix)
		return regexp.MustCompile(r)
	}
}

func matchAll(text string, regex *regexp.Regexp) ([]string, [][]int) {
	parsed := regex.FindAllString(text, -1)
	index := regex.FindAllStringIndex(text, -1)
	return parsed, index
}

func matchString(text string, regex *regexp.Regexp) []string {
	return regex.FindAllString(text, -1)
}

func matchIndex(text string, regex *regexp.Regexp) [][]int {
	return regex.FindAllStringIndex(text, -1)
}

func matchBetween(text string, offset int, patternStart *regexp.Regexp, patternEnd *regexp.Regexp, cutset string) []string {
	text = lygo_strings.Sub(text, offset, len(text))
	response := make([]string, 0)

	indexesStart := matchIndex(text, patternStart) // [][]int
	for _, indexStart := range indexesStart {
		is := indexStart[0]
		ie := indexStart[1] // end of first pattern
		if is < ie {
			sub := lygo_strings.Sub(text, ie, len(text))
			indexesEnd := matchIndex(sub, patternEnd) // [][]int
			if len(indexesEnd) > 0 {
				indexEnd := indexesEnd[0][0]
				sub = lygo_strings.Sub(sub, 0, indexEnd)
				if len(cutset) > 0 {
					sub = strings.Trim(sub, cutset)
				}
				response = append(response, sub)
			} else {
				if len(cutset) > 0 {
					sub = strings.Trim(sub, cutset)
				}
				response = append(response, sub)
			}
		}
	}

	return response
}

func index(text string, regex *regexp.Regexp, offset int) []int {
	var response []int
	if nil != regex && len(text) > 0 {
		if offset < 0 {
			offset = 0
		}

		// shrink text starting from offset
		text = lygo_strings.Sub(text, offset, len(text))

		// get regexp match
		indexes := matchIndex(text, regex)

		if len(indexes) > 0 {
			for _, index := range indexes {
				response = append(response, index[0]+offset)
			}
		}
	}
	return response
}

func indexLenPair(text string, regex *regexp.Regexp, offset int) [][]int {
	var response [][]int
	if nil != regex && len(text) > 0 {
		if offset < 0 {
			offset = 0
		}

		// shrink text starting from offset
		text = lygo_strings.Sub(text, offset, len(text))

		// get regexp match
		indexes := matchIndex(text, regex)

		if len(indexes) > 0 {
			for _, index := range indexes {
				pair := make([]int, 2)
				pair[0] = index[0] + offset
				pair[1] = index[1] - index[0]
				response = append(response, pair)
			}
		}
	}
	return response
}
