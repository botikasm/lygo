package lygo_ocr


type OcrParams struct {
	FileName  string
	Trim      bool
	Verbose   bool
	Whitelist string
}

// Create an empty initialized instance of OcrParams
func NewOcrParams() *OcrParams {
	response := new(OcrParams)
	response.Trim = false
	response.Whitelist = ""

	return response
}
