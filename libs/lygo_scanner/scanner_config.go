package lygo_scanner

import "encoding/json"

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type ScannerConfigArray []ScannerConfig

type ScannerConfig struct {
	Uid         string               `json:"uid"`
	Description string               `json:"description"`
	Rotate      bool                 `json:"rotate"`
	Sharpen     bool                 `json:"sharpen"`
	Areas       []*ScannerConfigArea `json:"areas"`

	// optional NLP
	Entities []*ScannerConfigEntity `json:"entities"`
}

type ScannerConfigArea struct {
	Uid    string `json:"uid"`
	X      int    `json:"X"`
	Y      int    `json:"Y"`
	Width  uint   `json:"width"`
	Height uint   `json:"height"`

	// optional NLP
	Entities []*ScannerConfigEntity `json:"entities"`
}

type ScannerConfigEntity struct {
	Uid         string   `json:"uid"`
	Description string   `json:"description"`
	Intent      string   `json:"intent"`
	Score       int      `json:"score"`
	Values      []string `json:"values"`
}

//----------------------------------------------------------------------------------------------------------------------
//	ScannerConfigArray
//----------------------------------------------------------------------------------------------------------------------

/*
@param text : Json Array of rules to parse
[
	{
		"uid": "mod_70",
		"description": "",
		"rotate": true,
		"sharpen": true,
		"one_document_per_page": true,
		"areas": [
		  {
			"uid": "doc_type",
			"description": "TEST AREA",
			"X": 0,
			"Y": 0,
			"width": 10000,
			"height": 10000
		  }
		]
  	}
]
*/
func (config *ScannerConfigArray) Parse(text string) error {
	return json.Unmarshal([]byte(text), &config)
}

func (config *ScannerConfigArray) ToString() string {
	b, err := json.Marshal(&config)
	if nil == err {
		return string(b)
	}
	return ""
}

func (config *ScannerConfigArray) Clone() *ScannerConfigArray {
	clone := new(ScannerConfigArray)
	for _, item := range config.Items() {
		clone.Add(&item)
	}
	return clone
}

func (config *ScannerConfigArray) Items() []ScannerConfig {
	return []ScannerConfig(*config)
}

func (config *ScannerConfigArray) Add(item *ScannerConfig) {
	v := []ScannerConfig(*config)
	*config = append(v, *item)
}

func (config *ScannerConfigArray) Remove(item *ScannerConfig) {
	v := []ScannerConfig(*config)
	var array []ScannerConfig
	for _, s := range v {
		if s.Uid != item.Uid {
			array = append(array, s)
		}
	}
	if nil == array {
		*config = make([]ScannerConfig, 0)
	} else {
		*config = array
	}
}

func (config *ScannerConfigArray) Contains(item *ScannerConfig) bool {
	v := []ScannerConfig(*config)
	for _, s := range v {
		if s.Uid == item.Uid {
			return true
		}
	}
	return false
}

func (config *ScannerConfigArray) AddAll(addItems []ScannerConfig) {
	v := []ScannerConfig(*config)
	items := make([]ScannerConfig, 0)
	for _, item := range addItems {
		if nil != &item && !config.Contains(&item) {
			items = append(items, item)
		}
	}
	*config = append(v, items...)
}

func (config *ScannerConfigArray) Count() int {
	v := []ScannerConfig(*config)
	return len(v)
}

func (config *ScannerConfigArray) IsEmpty() bool {
	return config.Count() == 0
}

func (config *ScannerConfigArray) Get(name string) *ScannerConfig {
	for _, item := range *config {
		if item.Uid == name {
			return &item
		}
	}
	return nil
}

//----------------------------------------------------------------------------------------------------------------------
//	ScannerConfig
//----------------------------------------------------------------------------------------------------------------------

func NewScannerConfig() *ScannerConfig {
	result := new(ScannerConfig)

	result.Rotate = false  // avoid image rotation
	result.Sharpen = false // avoid image sharpen

	result.Areas = nil

	return result
}

/*
	{
		"uid": "mod_70",
		"description": "",
		"rotate": true,
		"sharpen": true,
		"areas": [
		  {
			"uid": "doc_type",
			"description": "TEST AREA",
			"X": 0,
			"Y": 0,
			"width": 10000,
			"height": 10000
		  }
		]
  	}
*/
func (instance *ScannerConfig) Parse(text string) error {
	return json.Unmarshal([]byte(text), &instance)
}

func (instance *ScannerConfig) ToString() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}

func (instance *ScannerConfig) AllEntities() []*ScannerConfigEntity {
	response := make([]*ScannerConfigEntity, 0)
	if nil!=instance.Entities{
		response = append(response, instance.Entities...)
	}
	if nil != instance.Areas {
		for _, area := range instance.Areas {
			response = append(response, area.Entities...)
		}
	}
	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	ScannerConfigArea
//----------------------------------------------------------------------------------------------------------------------

func (instance *ScannerConfigArea) Parse(text string) error {
	return json.Unmarshal([]byte(text), &instance)
}

func (instance *ScannerConfigArea) ToString() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}

func (instance *ScannerConfigArea) IsEnabled() bool {
	return instance.Width > 0 && instance.Height > 0
}

func (instance *ScannerConfigArea) IsFullPage() bool {
	return instance.Width == instance.Height && instance.Width == 10000
}

//----------------------------------------------------------------------------------------------------------------------
//	ScannerConfigEntity
//----------------------------------------------------------------------------------------------------------------------

func (instance *ScannerConfigEntity) Parse(text string) error {
	return json.Unmarshal([]byte(text), &instance)
}

func (instance *ScannerConfigEntity) ToString() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}
