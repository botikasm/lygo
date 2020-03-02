package lygo_scanner

const _MAX_ROTATE = 4

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type ScannerResponse struct {
	Uid       string
	Original  string
	Params    *ScannerConfigArray
	Tag       string
	Documents []*ScannerDocument
	ElapsedMs int
}

//----------------------------------------------------------------------------------------------------------------------
//	ScannerResponse
//----------------------------------------------------------------------------------------------------------------------

func (instance *ScannerResponse) BestDocument() *ScannerDocument {
	var response *ScannerDocument
	var score float32
	score = -1.0

	for _, document := range instance.Documents {
		if nil != document {
			for _, page := range document.Pages {
				job := page.BestJob()
				if nil != job {
					for _, item := range job.Areas {
						if nil != item && nil != item.Nlp {
							if item.Nlp.Score > score {
								score = item.Nlp.Score
								response = document
							}
						}
					}
				}
			}

		}
	}
	return response
}
