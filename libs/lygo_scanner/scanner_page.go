package lygo_scanner

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type ScannerPage struct {
	Id       int
	FileName string
	Jobs     [_MAX_ROTATE]*ScannerPageJob
}

type ScannerPageJob struct {
	Index    int // page index
	Parent   *ScannerPage
	FileName string
	Elapsed  int
	Error    error
	Areas    []*ScannerPageJobArea
}

type ScannerPageJobArea struct {
	Index       int // position index in areas array
	Parent      *ScannerPageJob
	Uid         string
	FileName    string
	Text        string
	IsFullPage  bool
	Coordinates *ScannerPageJobAreaCoordinates
	Error       error
	Nlp         *ScannerPageJobAreaNlpResponse
}

type ScannerPageJobAreaCoordinates struct {
	X      int
	Y      int
	Width  uint
	Height uint
}

type ScannerPageJobAreaNlpResponse struct {
	Parent          *ScannerPageJobArea
	Score           float32
	Elapsed         int
	Entities        map[string][]interface{}
	IntentEntityUid string
	IntentUid       string
}

//----------------------------------------------------------------------------------------------------------------------
//	ScannerPage
//----------------------------------------------------------------------------------------------------------------------

func (instance *ScannerPage) HasIntent() bool {
	for _, job := range instance.Jobs {
		if nil != job {
			if nil != job.Areas {
				for _, item := range job.Areas {
					if nil != item && nil != item.Nlp {
						if len(item.Nlp.IntentUid) > 0 {
							return true
						}
					}
				}
			}
		}
	}
	return false
}

func (instance *ScannerPage) IntentUid() string {
	for _, job := range instance.Jobs {
		if nil != job {
			if nil != job.Areas {
				for _, item := range job.Areas {
					if nil != item && nil != item.Nlp {
						if len(item.Nlp.IntentUid) > 0 {
							return item.Nlp.IntentUid
						}
					}
				}
			}
		}
	}
	return ""
}

// Return best scanner elaboration (nlp matching with image transformation) in a page
func (instance *ScannerPage) BestJob() *ScannerPageJob {
	var response *ScannerPageJob
	var score float32
	score = -1.0
	for _, job := range instance.Jobs {
		if nil != job {
			if nil != job.Areas {
				for _, item := range job.Areas {
					if nil != item && nil != item.Nlp {
						if item.Nlp.Score > score {
							score = item.Nlp.Score
							response = job
						}
					}
				}
			}
		}
	}
	return response
}

func (instance *ScannerPage) Entities() map[string][]interface{} {
	var response map[string][]interface{}
	response = make(map[string][]interface{})
	job := instance.BestJob()
	if nil != job {
		for _, item := range job.Areas {
			if nil != item && nil != item.Nlp {
				if item.Nlp.Score > 0 {
					// copy values
					for k, v := range item.Nlp.Entities {
						response[k] = append(response[k], v...)
					}
				}
			}
		}
	}

	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	ScannerPageJob
//----------------------------------------------------------------------------------------------------------------------

func (instance *ScannerPageJob) Entities(minScore float32) map[string][]interface{} {
	var response map[string][]interface{}
	response = make(map[string][]interface{})
	if len(instance.Areas) > 0 {
		for _, item := range instance.Areas {
			if nil != item && nil != item.Nlp {
				if item.Nlp.Score >= minScore {
					// copy values
					for k, v := range item.Nlp.Entities {
						response[k] = append(response[k], v...)
					}
				}
			}
		}
	}
	return response
}
