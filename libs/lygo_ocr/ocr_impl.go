package lygo_ocr

import (
	"errors"
	"github.com/otiai10/gosseract"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func Version() string {
	return gosseract.Version()
}

func ReadText(p *OcrParams) (string, error) {
	if len(p.FileName) == 0 {
		return "", errors.New("missing 'FileName' parameter")
	}

	client := gosseract.NewClient()
	defer client.Close()

	client.Trim = p.Trim
	if !p.Verbose{
		client.DisableOutput()
	}

	client.SetWhitelist(p.Whitelist)

	// read all text
	client.SetImage(p.FileName)
	text, err := client.Text()

	return text, err
}
