package lygo_nlpdetect

import (
	"encoding/json"
	"log"
	"os"
	"regexp"
	"sort"
	"strings"
	"unicode/utf8"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t
//----------------------------------------------------------------------------------------------------------------------

const MAX_DIFFERENCE = float64(300)
const MAX_LENGTH = 2048
const MIN_LENGTH = 3

//----------------------------------------------------------------------------------------------------------------------
//	f i e l d s
//----------------------------------------------------------------------------------------------------------------------

var scripts map[string]languages
var expressions map[string]regexp.Regexp

//----------------------------------------------------------------------------------------------------------------------
//	i n i t
//----------------------------------------------------------------------------------------------------------------------

//Init `scripts` and `expressions` dictionaries
func init() {
	var scriptData map[string]interface{}
	err := json.Unmarshal([]byte(SCRIPT_DATA), &scriptData)
	if err != nil {
		log.Printf("Error during languages decoding: %v\n", err)
		os.Exit(1)
	}

	scripts = make(map[string]languages)
	for scriptCode, scriptValue := range scriptData {
		languages := make(map[string]trigrams)

		lang := scriptValue.(map[string]interface{})
		for code, trigramsRaw := range lang {
			trigramsString := trigramsRaw.(string)
			trigramsStringArray := strings.Split(trigramsString, "|")
			trigrams := make(map[string]int)
			for i := len(trigramsStringArray) - 1; i >= 0; i-- {
				trigrams[trigramsStringArray[i]] = i + 1
			}
			languages[code] = trigrams
		}

		scripts[scriptCode] = languages
	}

	var expressionsData map[string]interface{}
	err = json.Unmarshal([]byte(EXPRESSION_DATA), &expressionsData)
	if err != nil {
		log.Printf("Error during expressions decoding: %v\n", err)
		os.Exit(1)
	}

	expressions = make(map[string]regexp.Regexp)
	for code, v := range expressionsData {
		expressions[code] = *regexp.MustCompile(v.(string))
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------


//Get a list of probable languages the given value is written in.
func DetectWithFilters(value string, whitelist []string, blacklist []string) tuples {
	if len(value) < MIN_LENGTH {
		return singleLanguageTuples("und")
	}

	if len(value) > MAX_LENGTH {
		value = value[0:MAX_LENGTH]
	}

	code := getTopScript(value)

	if _, ok := scripts[code]; !ok {
		return singleLanguageTuples(code)
	}

	return normalize(value, getDistances(getTrigramsAsTuples(value), scripts[code], whitelist, blacklist))
}

//Get a list of probable languages the given value is written in.
func Detect(value string) tuples {
	return DetectWithFilters(value, make([]string, 0), make([]string, 0))
}

//Get the most probable language for the given value.
func DetectOne(value string) tuple {
	return Detect(value)[0]
}

//Get the most probable language for the given value.
func DetectOneWithFilters(value string, whitelist []string, blacklist []string) tuple {
	return DetectWithFilters(value, whitelist, blacklist)[0]
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

//Filter `languages` by removing languages in `blacklist`, or including languages in `whitelist`.
func filterLanguages(langs languages, whitelist []string, blacklist []string) languages {
	filteredLanguages := make(languages)

	if len(whitelist) == 0 && len(blacklist) == 0 {
		return langs
	}

	if len(whitelist) == 0 {
		filteredLanguages = langs
	} else {
		for _, code := range whitelist {
			filteredLanguages[code] = langs[code]
		}
	}

	for _, code := range blacklist {
		if _, ok := filteredLanguages[code]; ok {
			delete(filteredLanguages, code)
		}
	}
	return filteredLanguages
}

//Get the distance between `trigrams`, and a trigram `model`.
func getDistance(trigrams tuples, model trigrams) float64 {
	distance := float64(0)

	for _, t := range trigrams {
		difference := float64(0)
		if modelVal, ok := model[t.Code]; ok {
			difference = t.Count - float64(modelVal)

			if difference < 0 {
				difference = -difference
			}
		} else {
			difference = MAX_DIFFERENCE
		}

		distance += difference
	}

	return distance
}

//Get the distance between `trigrams`, and multiple trigram dictionaries `languages`.
func getDistances(trigrams tuples, languages languages, whitelist []string, blacklist []string) tuples {
	filteredLanguages := filterLanguages(languages, whitelist, blacklist)
	tuples := make(tuples, 0)

	for code, language := range filteredLanguages {
		dis := getDistance(trigrams, language)
		t := tuple{Code: code, Count: dis}
		tuples = append(tuples, t)
	}

	return tuples
}

//Get the occurrence ratio of `expression` for `value`.
func getOccurrence(value string, expression regexp.Regexp) float64 {
	count := len(expression.FindAllString(value, -1))

	if count < 1 {
		count = 0
	}

	return float64(count) / float64(len(value))
}

//From `scripts`, get the most occurring script for `value`
func getTopScript(value string) string {

	topCount := float64(-1)
	expressionCode := ""

	for code, e := range expressions {
		count := getOccurrence(value, e)

		if count > topCount {
			topCount = count
			expressionCode = code
		}
	}

	return expressionCode
}

//Create a single tuple as a list of tuples from a given language code.
func singleLanguageTuples(code string) tuples {
	tuples := make(tuples, 1)
	tuples[0] = tuple{Code: code, Count: 1}
	return tuples
}


func normalize(value string, distances tuples) tuples {
	sort.Sort(distances)

	min := distances[0].Count
	max := (float64(len(value)) * MAX_DIFFERENCE) - min

	for i, d := range distances {
		distances[i].Count = float64(1) - ((d.Count - min) / max)
	}

	return distances
}

func clean(value string) string {
	re := regexp.MustCompile("[\u0021-\u0040]+")
	value = re.ReplaceAllString(value, " ")

	re = regexp.MustCompile(`\s+`)
	value = re.ReplaceAllString(value, " ")

	value = strings.ToLower(value)
	value = strings.Trim(value, " ")

	return value
}

func getTrigrams(value string) []string {
	value = " " + clean(value) + " "
	res := make([]string, 0)
	val := NewString(value) // utf8string.NewString(value)
	i := 0
	for i+3 < utf8.RuneCountInString(value) {
		res = append(res, val.Slice(i, i+3))
		i++
	}

	return res
}

func getTrigramsAsMap(value string) map[string]int {
	trigrams := getTrigrams(value)
	res := make(map[string]int)
	for _, t := range trigrams {
		if _, ok := res[t]; !ok {
			res[t] = 0
		}
		res[t]++
	}

	return res
}

func getTrigramsAsTuples(value string) []tuple {
	trigrams := getTrigramsAsMap(value)
	res := make([]tuple, len(trigrams))
	i := 0
	for code, count := range trigrams {
		res[i] = tuple{Code: code, Count: float64(count)}
		i++
	}

	return res
}
