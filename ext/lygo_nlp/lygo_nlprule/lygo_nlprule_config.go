package lygo_nlprule

import "encoding/json"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type NlpRuleConfigArray []NlpRuleConfigIntent

type NlpRuleConfigIntent struct {
	Uid         string                `json:"uid"`
	Description string                `json:"description"`
	Entities    []NlpRuleConfigEntity `json:"entities"`
}

type NlpRuleConfigEntity struct {
	Uid         string   `json:"uid"`
	Description string   `json:"description"`
	Intent      string   `json:"intent"`
	Score       int      `json:"score"`
	Values      []string `json:"values"`
}

//----------------------------------------------------------------------------------------------------------------------
//	NlpRuleConfigArray
//----------------------------------------------------------------------------------------------------------------------

/*
@param text : Json Array of rules to parse
[
	{
		"uid": "mod_80",
		"description": "Intent or document identifier fo Cod. 80",
		"entities": [
			{
				"uid": "type",
				"description": "document type. Lookup for 'Cod. 80'",
				"intent": "",
				"score": 100,
				"values": ["$regexps.HasMatch('?od??80') || $regexps.HasMatch('?od?80')"]
			}
		]
	}
]
*/
func (config *NlpRuleConfigArray) Parse(text string) error {
	return json.Unmarshal([]byte(text), &config)
}

func (config *NlpRuleConfigArray) ToString() string {
	b, err := json.Marshal(&config)
	if nil == err {
		return string(b)
	}
	return ""
}

func (config *NlpRuleConfigArray) Items() []NlpRuleConfigIntent {
	return []NlpRuleConfigIntent(*config)
}

func (config *NlpRuleConfigArray) Add(item *NlpRuleConfigIntent) {
	v := []NlpRuleConfigIntent(*config)
	*config = append(v, *item)
}

func (config *NlpRuleConfigArray) AddAll(items []NlpRuleConfigIntent) {
	v := []NlpRuleConfigIntent(*config)
	*config = append(v, items...)
}

func (config *NlpRuleConfigArray) Count() int {
	v := []NlpRuleConfigIntent(*config)
	return len(v)
}

func (config *NlpRuleConfigArray) IsEmpty() bool {
	return config.Count() == 0
}

func (config *NlpRuleConfigArray) Get(name string) *NlpRuleConfigIntent {
	for _, item := range *config {
		if item.Uid == name {
			return &item
		}
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	NlpRuleConfigIntent
//----------------------------------------------------------------------------------------------------------------------

func (instance *NlpRuleConfigIntent) Parse(text string) error {
	return json.Unmarshal([]byte(text), &instance)
}

func (instance *NlpRuleConfigIntent) ToString() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}

//----------------------------------------------------------------------------------------------------------------------
//	NlpRuleConfigEntity
//----------------------------------------------------------------------------------------------------------------------

func (instance *NlpRuleConfigEntity) Parse(text string) error {
	return json.Unmarshal([]byte(text), &instance)
}

func (instance *NlpRuleConfigEntity) ToString() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}
