package lygo_nlprule

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type NlpRuleEngineResponse struct {
	ElapsedMs int
	Items     []NlpRuleEngineResponseItem
}

type NlpRuleEngineResponseItem struct {
	Score        float32
	IntentUid    string
	IntentScore  float32
	IntentEntity string // name of entity that is also an intent
	Entities     []NlpRuleEngineResponseItemEntity
}
type NlpRuleEngineResponseItemEntity struct {
	Uid    string
	Rules  []string
	Errors []string
	Values []interface{}
}

//----------------------------------------------------------------------------------------------------------------------
//	NlpRuleEngineResponse
//----------------------------------------------------------------------------------------------------------------------

func (instance *NlpRuleEngineResponse) IntentEntityUid() string {
	if nil != instance.Items {
		for _, item := range instance.Items {
			if len(item.IntentEntity) > 0 {
				return item.IntentEntity
			}
		}
	}
	return ""
}

func (instance *NlpRuleEngineResponse) Values() map[string][]interface{} {
	response := make(map[string][]interface{})
	if nil != instance.Items {
		for _, item := range instance.Items {
			if nil != item.Entities {
				for _, entity := range item.Entities {
					key := entity.Uid
					values := entity.Values
					if len(key) > 0 {
						response[key] = values
					}
				}
			}
		}
	}
	return response
}
