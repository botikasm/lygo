package lygo_nlprule

import (
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_stopwatch"
	"github.com/botikasm/lygo/base/lygo_strings"
	"github.com/botikasm/lygo/ext/lygo_scripting"
	"github.com/dop251/goja"
	"path"
	"sort"
	"strings"
)

const VARIABLE_PREFIX = "VAR_"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type NlpRuleEngine struct {
	Config *NlpRuleConfigArray
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewRuleEngine(config interface{}) *NlpRuleEngine {
	response := new(NlpRuleEngine)

	if isString, s := lygo_conv.IsString(config); isString {
		// STRING (json text or filename)
		c := new(NlpRuleConfigArray)
		if strings.HasPrefix(s, "[") {
			// directly parse as json text
			_ = c.Parse(s)
		} else {
			if len(path.Base(s)) > 0 {
				// read text from file
				text, err := lygo_io.ReadTextFromFile(s)
				if nil == err {
					// parse as json text
					_ = c.Parse(text)
				}
			}
		}
		response.Config = c
	} else {
		// STRUCT
		v := config.(*NlpRuleConfigArray)
		if nil != v {
			response.Config = v
		}
	}

	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (engine *NlpRuleEngine) HasConfig() bool {
	return nil != engine.Config && !engine.Config.IsEmpty()
}

// Evaluate all rules and return a response containing all reached items sorted descending by Score. Best score
// is on the top of list
// @text: Test to evaluate
// @context: Optional context parameters useful for evaluation
// @limitScore: If greater than zero, engine stop evaluating when score is equal or greater then limitScore
// Returns an array sorted desc from best score to worst
func (engine *NlpRuleEngine) Eval(text string, context map[string]interface{}, limitScore float32) *NlpRuleEngineResponse {
	stopwatch := lygo_stopwatch.New()
	stopwatch.Start()

	response := new(NlpRuleEngineResponse)
	response.ElapsedMs = 0

	runtime := engine.initExpressionEngine(text, context)

	// LOOP on all configuration items
	items := engine.Config.Items()
	for _, item := range items {
		if len(item.Uid) > 0 {
			itemResponse := engine.eval(runtime, &item)
			response.Items = append(response.Items, *itemResponse)

			// check LIMIT SCORE
			if len(response.Items) > 0 {
				last := response.Items[len(response.Items)-1]
				if limitScore > 0 && last.Score >= limitScore {
					break // EXIT if reached a limit score
				}
			}
		}
	}

	stopwatch.Stop()
	response.ElapsedMs = stopwatch.Milliseconds()

	// sort by score DESC
	sort.Slice(response.Items, func(i, j int) bool {
		return response.Items[i].Score > response.Items[j].Score
	})

	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (engine *NlpRuleEngine) initExpressionEngine(text string, context map[string]interface{}) *lygo_scripting.ScripEngine {
	runtime := lygo_scripting.New()

	// init global context as text
	runtime.SetToolsContext(text)

	// add text also as variable
	runtime.Set(VARIABLE_PREFIX+"text", text)

	if nil != context {
		for k, v := range context {
			if nil != v && nil == runtime.Get(k) {
				runtime.Set(VARIABLE_PREFIX+k, v)
			}
		}
	}

	return runtime
}

// Evaluate a rule node.
/*
	{
    "uid": "mod_80",
    "description": "Intent or document identifier fo Cod. 80",
    "entities": [
      {
        "uid": "doc_type",
        "description": "[FIRST IS AN INTENT TOO] document type. Lookup for 'Cod. 80'",
        "intent": "mod_80",
        "score": 100,
        "values": [
          "$regexps.MatchFirst('?od??80')"
        ]
      },
      {
        "uid": "price",
        "description": "price to pay",
        "score": 100,
        "values": []
      }
    ]
  }
*/
func (engine *NlpRuleEngine) eval(runtime *lygo_scripting.ScripEngine, item *NlpRuleConfigIntent) *NlpRuleEngineResponseItem {
	response := new(NlpRuleEngineResponseItem)
	response.Score = 0
	response.IntentScore = 0
	response.IntentUid = ""
	response.IntentEntity = ""

	entityRules := item.Entities
	for _, rule := range entityRules {
		// mark as intent
		if engine.isIntent(&rule) {
			response.IntentEntity = rule.Uid
		}

		// set script runtime references for console log files
		runtime.Root = lygo_paths.Concat("logging", item.Uid)
		runtime.Name = rule.Uid

		// solve all expressions find in rule
		entities := engine.solveRule(runtime, rule)

		if len(entities.Values) > 0 {

			// append entities
			response.Entities = append(response.Entities, *entities)

			if engine.isIntent(&rule) {
				// found an INTENT with values
				response.IntentUid = item.Uid
				response.IntentScore = float32(rule.Score * len(entities.Values))
			}

			// update global score
			if response.IntentScore > 0 {
				response.Score = response.Score + response.IntentScore
			} else {
				response.Score = response.Score + float32(rule.Score*len(entities.Values))
			}
		} else {
			if engine.isIntent(&rule) {
				// found an intent with NO values
				// EXIT loop and does not evaluate entities
				break
			}
		}
	}

	return response
}

func (engine *NlpRuleEngine) isIntent(rule *NlpRuleConfigEntity) bool {
	return len(rule.Intent) > 0
}

func (engine *NlpRuleEngine) solveRule(runtime *lygo_scripting.ScripEngine, rule NlpRuleConfigEntity) *NlpRuleEngineResponseItemEntity {
	response := new(NlpRuleEngineResponseItemEntity)
	response.Uid = rule.Uid

	scripts := engine.solveScriptPath(rule.Values)
	// solve values
	for _, expr := range scripts {
		if len(expr) > 0 {

			response.Rules = append(response.Rules, expr) // current expression

			// solve
			v, err := runtime.RunString(expr)

			// check error or value
			if nil == err {
				response.Values = append(response.Values, convert(v)...)
			} else {
				response.Errors = append(response.Errors, err.Error())
			}
		}
	}

	return response
}

func (document *NlpRuleEngine) solveScriptPath(values []string) []string {
	response := make([]string, 0)
	for _, value := range values {
		if len(value) > 0 {
			if strings.HasPrefix(value, "file://") {
				path := strings.Replace(value, "file://", "", 1)
				path = lygo_paths.WorkspacePath(path)
				if b, _ := lygo_paths.Exists(path); b {
					s, err := lygo_io.ReadTextFromFile(path)
					if len(s) > 0 {
						response = append(response, s)
					} else if nil != err {
						response = append(response, lygo_strings.Format("(function(){return '%s'})()", err.Error()))
					}
				}
			} else {
				response = append(response, value)
			}
		}
	}
	return response
}

func convert(raw_value goja.Value) []interface{} {
	response := make([]interface{}, 0)
	if nil != raw_value {
		value := raw_value.Export()
		if b, _ := lygo_conv.IsArray(value); b {
			a := toArray(value)
			if len(a) > 0 {
				response = append(response, a...)
			}
		} else {
			if isValidValue(value) {
				response = append(response, value)
			}
		}
	}

	return response
}

func toArray(array interface{}) []interface{} {
	source := lygo_conv.ToArray(array)
	var response []interface{}
	for _, v := range source {
		if isValidValue(v) {
			response = append(response, v)
		}
	}
	return response
}

func isValidValue(ivalue interface{}) bool {
	value := lygo_conv.ToString(ivalue)
	value = strings.Trim(value, " ")
	return len(value) > 0 && value != "false"
}
