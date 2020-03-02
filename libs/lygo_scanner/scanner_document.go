package lygo_scanner

import (
	"errors"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/botikasm/lygo/base/lygo_rnd"
	"github.com/botikasm/lygo/base/lygo_stopwatch"
	"github.com/botikasm/lygo/base/lygo_strings"
	"github.com/botikasm/lygo/ext/lygo_nlp/lygo_nlprule"
	"github.com/botikasm/lygo/libs/lygo_images"
	"github.com/botikasm/lygo/libs/lygo_ocr"
	"path/filepath"
	"strings"
	"sync"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type ScannerDocument struct {
	Root      string
	Id        string
	Workspace string
	Pages     []*ScannerPage // Jobs      [_MAX_ROTATE]*ScannerPageJob
	Params    *ScannerConfig
	ModelUid  string // assigned if found a matching
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewScannerDocument(root string, pages []string, params ScannerConfig) *ScannerDocument {

	response := NewScannerDocumentEmpty(root, pages, params)

	// init workspace
	lygo_paths.Mkdir(response.Workspace)

	// async execution in goroutines
	response.initialize(false)

	return response
}

func NewScannerDocumentEmpty(root string, pages []string, params ScannerConfig) *ScannerDocument {
	id := lygo_rnd.UuidDefault("scannerdoc_default")
	workspace := filepath.Join(root, "documents", id)

	response := &ScannerDocument{
		Root:      root,
		Id:        id,
		Workspace: workspace,
		Params:    &params,
	}

	// creates pages
	for index, filePage := range pages {
		page := new(ScannerPage)
		page.Id = index
		page.FileName = filePage
		response.Pages = append(response.Pages, page)
	}

	return response
}

func NewScannerDocumentDebug(root string, pages []string, params ScannerConfig) *ScannerDocument {

	response := NewScannerDocumentEmpty(root, pages, params)

	// init workspace
	lygo_paths.Mkdir(response.Workspace)

	// run just a debut session.
	response.initialize(true)

	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (document *ScannerDocument) Uid() string {
	return document.Params.Uid
}

func (document *ScannerDocument) BestPage() *ScannerPage {
	var response *ScannerPage
	var score float32
	for _, page := range document.Pages {
		for _, job := range page.Jobs {
			if nil != job {
				for _, item := range job.Areas {
					if nil != item && nil != item.Nlp {
						if item.Nlp.Score > score {
							score = item.Nlp.Score
							response = page
						}
					}
				}
			}
		}
	}
	return response
}

func (document *ScannerDocument) BestJob() *ScannerPageJob {
	var response *ScannerPageJob
	var score float32
	for _, page := range document.Pages {
		for _, job := range page.Jobs {
			if nil != job {
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

func (document *ScannerDocument) BestJobArea() *ScannerPageJobArea {
	var response *ScannerPageJobArea
	var score float32
	for _, page := range document.Pages {
		for _, job := range page.Jobs {
			if nil != job {
				for _, item := range job.Areas {
					if nil != item && nil != item.Nlp {
						if item.Nlp.Score > score {
							score = item.Nlp.Score
							response = item
						}
					}
				}
			}
		}
	}
	return response
}

func (document *ScannerDocument) BestJobItemScore() float32 {
	best := document.BestJobArea()
	if nil != best && nil != best.Nlp {
		return best.Nlp.Score
	}
	return 0.0
}

func (document *ScannerDocument) BestEntities() map[string][]interface{} {
	var score float32
	var response map[string][]interface{}
	if nil != document {
		for _, page := range document.Pages {
			job := page.BestJob()
			if nil != job {
				for _, item := range job.Areas {
					if nil != item && nil != item.Nlp {
						if item.Nlp.Score > score {
							score = item.Nlp.Score
							response = item.Nlp.Entities
						}
					}
				}
			}
		}
	}
	return response
}

func (document *ScannerDocument) Entities(minScore float32) map[string][]interface{} {
	var response map[string][]interface{}
	response = make(map[string][]interface{})
	if nil != document {
		for _, page := range document.Pages {
			job := page.BestJob()
			if nil != job {
				entities := job.Entities(minScore)
				if nil != entities {
					for k, v := range entities {
						response[k] = append(response[k], v...)
					}
				}
			}
		}
	}
	return response
}

// Get total score of a document
func (document *ScannerDocument) Score() float32 {
	var response float32
	if nil != document {
		for _, page := range document.Pages {
			job := page.BestJob()
			if nil != job {
				for _, item := range job.Areas {
					if nil != item && nil != item.Nlp {
						response += item.Nlp.Score
					}
				}
			}
		}
	}
	return response
}

// Return all text found in pages.
// Text include also area text, so it's not exactly what expected from a page ocr
func (document *ScannerDocument) TextAll() string {
	var response string
	if nil != document {
		for _, page := range document.Pages {
			job := page.BestJob()
			if nil != job {
				for _, item := range job.Areas {
					if nil != item {
						if len(response) > 0 {
							response += "\n"
						}
						response += item.Text
					}
				}
			}
		}
	}
	return response
}

// Returns only Full Pages Text.
func (document *ScannerDocument) Text() string {
	var response string
	if nil != document {
		for pageIndex, page := range document.Pages {
			job := page.BestJob()
			if nil != job {
				for _, item := range job.Areas {
					if nil != item {
						if item.IsFullPage {
							newPageLine := lygo_strings.Format("--- page: %s ---\n", pageIndex+1)
							//if len(response) > 0 {
							response += newPageLine
							//}
							response += item.Text
						}
					}
				}
			}
		}
	}
	return response
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (document *ScannerDocument) initialize(debugMode bool) {
	if debugMode {
		count := len(document.Pages)
		if count > 0 {
			for _, page := range document.Pages {
				// read file name that is expected is a text file
				if "txt" == lygo_paths.ExtensionName(page.FileName) {
					text, _ := lygo_io.ReadTextFromFile(page.FileName)
					if len(text) > 0 {
						// creates job
						job := new(ScannerPageJob)
						job.Index = 0
						job.Parent = page
						job.FileName = page.FileName
						page.Jobs[0] = job
						// generate area item
						item := new(ScannerPageJobArea)
						item.Uid = ""
						item.FileName = page.FileName
						item.Text = text
						item.Nlp = document.nlp(page.Id, 0, item.Text, document.Params.AllEntities())
						item.Nlp.Parent = item

						// add area item
						job.Areas = append(job.Areas, item)
					}
				} else {
					// invalid file extension
				}
			}
		}
	} else {
		count := len(document.Pages)
		if count > 0 {

			// creates jobs rotate_group
			var group sync.WaitGroup
			for index, page := range document.Pages {
				if index > 0 {
					group.Add(1)
					go document.goScan(&group, page)
				} else {
					// first page is sync because resolve intent
					document.goRotateAndFillJob(page)
					intent := page.IntentUid()
					if len(intent) == 0 {
						break // avoid to scan all other pages
					} else {
						document.ModelUid = intent
					}
				}
			}
			// wait all jobs are done
			group.Wait()
		}
	}
}

func (document *ScannerDocument) goScan(group *sync.WaitGroup, page *ScannerPage) {
	defer group.Done()

	document.goRotateAndFillJob(page)
}

func (document *ScannerDocument) goRotateAndFillJob(page *ScannerPage) {
	// creates jobs rotate_group
	var rotate_group sync.WaitGroup

	if document.Params.Rotate {
		rotate_group.Add(_MAX_ROTATE)
		for i := 0; i < _MAX_ROTATE; i++ {
			go document.rotateAndFillJobs(&rotate_group, i, page)
		}
	} else {
		rotate_group.Add(1)
		go document.rotateAndFillJobs(&rotate_group, 0, page)
	}
	// wait all jobs are done
	rotate_group.Wait()
}

func (document *ScannerDocument) rotateAndFillJobs(group *sync.WaitGroup, index int, page *ScannerPage) {
	defer group.Done()

	stopwatch := lygo_stopwatch.New()
	stopwatch.Start()

	// creates job
	job := new(ScannerPageJob)
	job.Index = index
	job.Parent = page
	page.Jobs[index] = job

	rp := lygo_images.NewImageRotateParams()
	rp.Degree = float64(index * 90.0)
	rp.Source = page.FileName
	rp.Target = filepath.Join(document.Workspace, lygo_conv.ToString(index)+"_"+filepath.Base(page.FileName))
	// rotate
	job.Error = lygo_images.Rotate(rp)
	job.FileName = rp.Target

	if nil == job.Error {
		// read text
		job.Areas = append(job.Areas, document.scanItem(job, job.FileName)...)
	}

	stopwatch.Stop()

	job.Elapsed = stopwatch.Milliseconds()
}

func (document *ScannerDocument) scanItem(parent *ScannerPageJob, filename string) []*ScannerPageJobArea {
	var response []*ScannerPageJobArea

	areas := document.Params.Areas
	if nil == areas {
		// entire file
		item := document.readItem(filename)
		item.Parent = parent
		item.IsFullPage = true
		item.Coordinates = &ScannerPageJobAreaCoordinates{}

		// check if NLP is enabled
		if nil == item.Error && nil != document.Params.Entities {
			item.Nlp = document.nlp(parent.Parent.Id, 0, item.Text, document.Params.Entities)
			item.Nlp.Parent = item
		}
		response = append(response, item)
	} else {
		// area
		for index, area := range areas {
			var item *ScannerPageJobArea
			if nil != area {
				// crop
				newFilename := lygo_paths.ChangeFileNameWithSuffix(filename, lygo_strings.Format("-%s", area.Uid))
				cropParams := lygo_images.NewImageCropParams()
				cropParams.Source = filename
				cropParams.Target = newFilename
				cropParams.X = area.X
				cropParams.Y = area.Y
				cropParams.Width = area.Width
				cropParams.Height = area.Height
				cropError := lygo_images.Crop(cropParams)
				if nil == cropError {
					item = document.readItem(newFilename)
					item.Parent = parent
					item.Index = index
					item.IsFullPage = area.IsFullPage()
					item.Coordinates = &ScannerPageJobAreaCoordinates{
						X:      area.X,
						Y:      area.Y,
						Width:  area.Width,
						Height: area.Height,
					}
					// check if NLP is enabled
					if nil == item.Error && nil != area.Entities {
						item.Nlp = document.nlp(parent.Parent.Id, index, item.Text, area.Entities)
						item.Nlp.Parent = item
					}
				} else {
					item = new(ScannerPageJobArea)
					item.Parent = parent
					item.Uid = ""
					item.FileName = newFilename
					item.Error = cropError
				}
				response = append(response, item)

				// break if is entity validator and has no matching
				// Is no useful match other areas if main intent is not detected
				isFirstAreaWithIntent := index == 0 && nil != item.Nlp && len(item.Nlp.IntentEntityUid) > 0
				if isFirstAreaWithIntent {
					if item.Nlp.Score == 0 && len(document.ModelUid) == 0 {
						// no more area scan
						break
					}
				}
			} else {
				item = new(ScannerPageJobArea)
				item.Parent = parent
				item.Uid = ""
				item.FileName = filename
				item.Error = errors.New("MalformedArea")
			}
		}
	}

	return response
}

func (document *ScannerDocument) readItem(filename string) *ScannerPageJobArea {
	// creates empty item
	item := new(ScannerPageJobArea)
	item.Uid = ""
	item.FileName = filename

	// read text
	item.Text, item.Error = document.read(filename)
	if nil == item.Error {
		// generate text file
		textFile := lygo_paths.ChangeFileNameExtension(filename, ".txt")
		lygo_io.WriteTextToFile(item.Text, textFile)
	}
	return item
}

func (document *ScannerDocument) read(filename string) (string, error) {
	ocrParams := lygo_ocr.NewOcrParams()
	ocrParams.FileName = filename
	return lygo_ocr.ReadText(ocrParams)
}

// Add NLP scoring at job item
func (document *ScannerDocument) nlp(pageIndex int, areaIndex int, text string, entities []*ScannerConfigEntity) *ScannerPageJobAreaNlpResponse {

	// initialize response
	response := new(ScannerPageJobAreaNlpResponse)
	response.Entities = make(map[string][]interface{})

	// Build NLP configuration
	config := new(lygo_nlprule.NlpRuleConfigIntent)
	config.Uid = document.Params.Uid
	config.Description = document.Params.Description
	for _, entity := range entities {
		configEntity := new(lygo_nlprule.NlpRuleConfigEntity)
		configEntity.Parse(entity.ToString())
		configEntity.Values = document.solveScriptPath(configEntity.Values)
		config.Entities = append(config.Entities, *configEntity)
		// add empty to have all declared values
		response.Entities[entity.Uid] = make([]interface{}, 0)
	}

	// Build NLP engine for current AREA item
	engineConfig := new(lygo_nlprule.NlpRuleConfigArray)
	engineConfig.Add(config)
	engine := lygo_nlprule.NewRuleEngine(engineConfig)
	if engine.HasConfig() {

		// creates context for script execution
		context := make(map[string]interface{})
		context["page_index"] = pageIndex
		context["area_index"] = areaIndex

		// resolve expressions and get a response
		engineResponse := engine.Eval(text, context, 0) // eval all
		if nil != engineResponse && len(engineResponse.Items) > 0 {
			// should check if intent has been detected
			detected := engineResponse.Items[0] // get best score
			isIntent := document.Params.Uid == detected.IntentUid
			valid := detected.Score > 0
			response.IntentEntityUid = engineResponse.IntentEntityUid()
			if valid {
				response.Score = detected.Score
				response.Elapsed = engineResponse.ElapsedMs
				// response.Entities = engineResponse.Values()
				mapValues := engineResponse.Values()
				if nil != mapValues {
					for k, v := range mapValues {
						if nil != v {
							response.Entities[k] = v
						}
					}
				}
				if isIntent {
					response.IntentUid = detected.IntentUid
				}
			}
		}
	} else {
		// CONFIGURATION MISMATCH
		// TODO: handle this kind of configuration errors

	}

	return response
}

func (document *ScannerDocument) solveScriptPath(values []string) []string {
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
