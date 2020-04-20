package lygo_scanner

import (
	"errors"
	"fmt"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_rnd"
	"github.com/botikasm/lygo/base/lygo_stopwatch"
	"github.com/botikasm/lygo/base/lygo_strings"
	"github.com/botikasm/lygo/ext/lygo_logs"
	"github.com/botikasm/lygo/libs/lygo_images"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type Scanner struct {
	Id               string
	Workspace        string
	ConversionFormat string
	ConversionDpi    float64

	//-- p r i v a t e --//
	mutex sync.Mutex
}

type ScannerTask struct {
	Processed []*ScannerDocument
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewScanner(optWorkspace ...interface{}) *Scanner {
	id := lygo_rnd.UuidTimestamp() // get new timestamped UUID

	// init root
	tmpRoot := lygo_paths.GetTempRoot()
	if len(optWorkspace) == 1 {
		if b, v := lygo_conv.IsString(optWorkspace[0]); b {
			tmpRoot = v
		}
	}
	workspace := filepath.Join(tmpRoot, id)

	response := &Scanner{
		Id:               id,
		Workspace:        workspace,
		ConversionFormat: "jpg",
	}

	// init workspace
	lygo_paths.Mkdir(workspace)

	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (scanner *Scanner) SplitDocuments(fileName string, oneDocumentPerPage bool) ([][]string, error) {
	defer func() {
		if r := recover(); r != nil {
			// recovered from panic
			msg := lygo_strings.Format("Scanner.SplitDocuments ERROR: %s", r)
			lygo_logs.Error(msg)
			fmt.Println("Scanner.SplitDocuments ERROR: ", msg)
		}
	}()

	response := make([][]string, 0)

	// copy original to workspace
	original, err := scanner.copyToWorkspace(fileName)
	if nil == err {
		// try to convert file if is PDF
		var isConverted bool
		var pages []string
		isConverted, pages, err = scanner.convert(original)

		if nil == err {
			if !isConverted {
				var target string
				target, err = scanner.copyToPage(original)
				pages = append(pages, target)
			}

			if nil == err {

				if oneDocumentPerPage {
					// many documents: one page generate one document
					for _, page := range pages {
						response = append(response, []string{page})
					}
				} else {
					// single multi-page document
					response = append(response, pages)
				}
			}
		}
	}

	return response, err
}

// Read a file that can be a multi-pages single document or single-page documents, marching with all available
// configurations to detect best matching.
// In case of many documents with a single page, the output respect document order.
// This procedure test all documents with all configurations
func (scanner *Scanner) ReadDocuments(fileName string, oneDocumentPerPage bool,
	configArray *ScannerConfigArray) (*ScannerResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			// recovered from panic
			msg := lygo_strings.Format("Scanner.ReadDocuments ERROR: %s", r)
			lygo_logs.Error(msg)
			fmt.Println("Scanner.ReadDocuments ERROR: ", msg)
		}
	}()

	stopwatch := lygo_stopwatch.New()
	stopwatch.Start()

	workConfig := configArray.Clone()

	response := new(ScannerResponse)
	response.Uid = scanner.Id
	response.Original = fileName
	response.Params = configArray
	response.Documents = make([]*ScannerDocument, 0)
	response.ElapsedMs = 0

	documents, err := scanner.SplitDocuments(fileName, oneDocumentPerPage)
	if nil == err {

		for _, documentPages := range documents {
			var bestScore float32
			var bestDocument *ScannerDocument
			var emptyDocument *ScannerDocument
			var fallbackDocument *ScannerDocument
			emptyDocument = NewScannerDocumentEmpty(scanner.Workspace, documentPages, *new(ScannerConfig))

			// loop on all configurations to get the best match
			var group sync.WaitGroup
			task := new(ScannerTask)
			task.Processed = make([]*ScannerDocument, 0)

			for _, config := range workConfig.Items() {

				// async fill list
				group.Add(1)
				go scanner.goScan(&group, task, documentPages, config)
				/*
					document := NewScannerDocument(scanner.Workspace, documentPages, config)
					fallbackDocument = document
					if len(document.Pages) > 0 {
						bestItem := document.BestJobArea()
						if nil != bestItem && nil != bestItem.Nlp {
							if bestItem.Nlp.Score > bestScore {
								bestScore = bestItem.Nlp.Score
								bestDocument = document
							}
						}
					}*/
			}

			// wait elaboration
			group.Wait()

			// assign best document matching or fallback
			for _, document := range task.Processed {
				fallbackDocument = document
				if len(document.Pages) > 0 {
					bestItem := document.BestJobArea()
					if nil != bestItem && nil != bestItem.Nlp {
						if bestItem.Nlp.Score > bestScore {
							bestScore = bestItem.Nlp.Score
							bestDocument = document
						}
					}
				}
			}

			if nil != bestDocument {
				// assign a document to response
				response.Documents = append(response.Documents, bestDocument)
				// matched configuration cannot be removed
				// workConfig.Remove(bestDocument.Params)
			} else if nil != fallbackDocument {
				// not a best document, but something with parameters and data
				response.Documents = append(response.Documents, fallbackDocument)
			} else {
				// document not matching models passed in configuration "configArray"
				response.Documents = append(response.Documents, emptyDocument)
			}
		}
	}

	stopwatch.Stop()
	response.ElapsedMs = stopwatch.Milliseconds()

	return response, err
}

// Match a file that can be a multi-pages single document or single-page documents, with a single configuration.
// This is useful to detect which document best match with passed configuration
func (scanner *Scanner) MatchDocuments(fileName string, oneDocumentPerPage bool,
	config *ScannerConfig) (*ScannerResponse, error) {
	defer func() {
		if r := recover(); r != nil {
			// recovered from panic
			msg := lygo_strings.Format("Scanner.MatchDocuments ERROR: %s", r)
			lygo_logs.Error(msg)
			fmt.Println("Scanner.MatchDocuments ERROR: ", msg)
		}
	}()

	stopwatch := lygo_stopwatch.New()
	stopwatch.Start()

	params := new(ScannerConfigArray)
	params.Add(config)

	response := new(ScannerResponse)
	response.Params = params
	response.Documents = make([]*ScannerDocument, 0)
	response.ElapsedMs = 0

	documents, err := scanner.SplitDocuments(fileName, oneDocumentPerPage)
	if nil == err {
		for _, documentPages := range documents {
			document := NewScannerDocument(scanner.Workspace, documentPages, *config)
			response.Documents = append(response.Documents, document)
		}
	}

	stopwatch.Stop()
	response.ElapsedMs = stopwatch.Milliseconds()

	return response, err
}

func (scanner *Scanner) Remove() error {
	dir := scanner.Workspace
	return os.RemoveAll(dir)
}

func (scanner *Scanner) DebugTextDocument(documentPages []string,
	configArray *ScannerConfigArray) (*ScannerResponse, error) {

	stopwatch := lygo_stopwatch.New()
	stopwatch.Start()

	workConfig := configArray.Clone()

	response := new(ScannerResponse)
	response.Uid = scanner.Id
	response.Original = ""
	response.Params = configArray
	response.Documents = make([]*ScannerDocument, 0)
	response.ElapsedMs = 0

	var bestScore float32
	var bestDocument *ScannerDocument
	var emptyDocument *ScannerDocument
	var fallbackDocument *ScannerDocument
	emptyDocument = NewScannerDocumentEmpty(scanner.Workspace, documentPages, *new(ScannerConfig))

	// loop on all configurations to get the best match
	for _, config := range workConfig.Items() {
		document := NewScannerDocumentDebug(scanner.Workspace, documentPages, config)
		fallbackDocument = document
		if len(document.Pages) > 0 {
			bestItem := document.BestJobArea()
			if nil != bestItem && nil != bestItem.Nlp {
				if bestItem.Nlp.Score > bestScore {
					bestScore = bestItem.Nlp.Score
					bestDocument = document
				}
			}
		}
	}

	if nil != bestDocument {
		// assign a document to response
		response.Documents = append(response.Documents, bestDocument)
		// matched configuration cannot be removed
		// workConfig.Remove(bestDocument.Params)
	} else if nil != fallbackDocument {
		// not a best document, but something with parameters and data
		response.Documents = append(response.Documents, fallbackDocument)
	} else {
		// document not matching models passed in configuration "configArray"
		response.Documents = append(response.Documents, emptyDocument)
	}

	stopwatch.Stop()
	response.ElapsedMs = stopwatch.Milliseconds()

	return response, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (scanner *Scanner) copyToWorkspace(source string) (string, error) {
	target := filepath.Join(scanner.Workspace, filepath.Base(source))
	_, err := lygo_io.CopyFile(source, target)
	return target, err
}

func (scanner *Scanner) buildDocumentFileName(original string, newExtension string) string {
	base := filepath.Base(original)
	if len(newExtension) > 0 {
		base = strings.Replace(base, lygo_paths.Extension(base), newExtension, 1)
	}
	return filepath.Join(scanner.Workspace, "original", base)
}

func (scanner *Scanner) copyToPage(source string) (string, error) {
	name := lygo_paths.ChangeFileNameWithSuffix(filepath.Base(source), "-0")
	target := filepath.Join(scanner.Workspace, "original", name)
	_, err := lygo_io.CopyFile(source, target)
	return target, err
}

func (scanner *Scanner) convert(original string) (bool, []string, error) {
	path := lygo_paths.Absolute(original)
	ext := lygo_paths.Extension(path)
	if len(ext) == 0 {
		return false, nil, errors.New("missing file extension")
	}

	format := scanner.ConversionFormat
	dpi := scanner.ConversionDpi

	if strings.ToLower(ext) == ".pdf" {
		targetPath := scanner.buildDocumentFileName(path, "."+format)
		params := lygo_images.NewImageConvertParams()
		params.Source = path
		params.Target = targetPath
		params.Format = format
		params.YRes = dpi
		params.XRes = dpi

		pages, err := lygo_images.Convert(params)

		return true, pages, err
	} else {
		// not converted & no error
		return false, nil, nil
	}
}

func (scanner *Scanner) goScan(group *sync.WaitGroup, task *ScannerTask, documentPages []string, config ScannerConfig) {
	defer group.Done()

	workspace := scanner.Workspace

	// scan the document
	document := NewScannerDocument(workspace, documentPages, config)

	if nil != document {
		defer scanner.mutex.Unlock()
		scanner.mutex.Lock()
		task.Processed = append(task.Processed, document)
	}
}
