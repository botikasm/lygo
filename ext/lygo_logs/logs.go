package lygo_logs

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_paths"
	"github.com/sirupsen/logrus"
	"os"
)

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t a n t s
//----------------------------------------------------------------------------------------------------------------------

// Program Workspace (child) logging folder. Here are all logs
const ROOT string = "logging"

const (
	FORMATTER_TEXT = 0
	FORMATTER_JSON = 1
)

const (
	LEVEL_PANIC = iota
	LEVEL_ERROR
	LEVEL_WARN
	LEVEL_INFO
	LEVEL_DEBUG
	LEVEL_TRACE
)

const (
	OUTPUT_CONSOLE = iota
	OUTPUT_FILE
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type LogParams struct {
	Level int
	Args  []interface{}
}

type LogContext struct {
	Caller string
	Data   interface{}
}

//----------------------------------------------------------------------------------------------------------------------
//	f i e l d s
//----------------------------------------------------------------------------------------------------------------------

var _logger *logrus.Logger // main logger
var _initialized bool
var _root, _fileName string
var _formatter int
var _level int
var _output int
var _channel chan *LogParams

//----------------------------------------------------------------------------------------------------------------------
//	i n i t
//----------------------------------------------------------------------------------------------------------------------

func init() {
	_initialized = false

	_logger = logrus.New()
	_logger.SetOutput(os.Stdout)
	_logger.SetLevel(logrus.InfoLevel)
	_logger.SetReportCaller(false)

	SetFormatter(FORMATTER_TEXT)
	SetLevel(LEVEL_INFO)
	SetOutput(OUTPUT_CONSOLE)

}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func GetRoot() string {
	return _root
}

func GetFileName() string {
	return _fileName
}

func SetFormatter(value int) {
	_formatter = value
	switch _formatter {
	case FORMATTER_TEXT:
		_logger.SetFormatter(&logrus.TextFormatter{})
	case FORMATTER_JSON:
		_logger.SetFormatter(&logrus.JSONFormatter{})
	default:
		_logger.SetFormatter(&logrus.TextFormatter{})
	}
}

func SetLevel(value int) {
	_level = value
	switch _level {
	case LEVEL_PANIC:
		_logger.SetLevel(logrus.PanicLevel)
	case LEVEL_ERROR:
		_logger.SetLevel(logrus.ErrorLevel)
	case LEVEL_WARN:
		_logger.SetLevel(logrus.WarnLevel)
	case LEVEL_INFO:
		_logger.SetLevel(logrus.InfoLevel)
	case LEVEL_DEBUG:
		_logger.SetLevel(logrus.DebugLevel)
	case LEVEL_TRACE:
		_logger.SetLevel(logrus.TraceLevel)
	default:
		_logger.SetLevel(logrus.WarnLevel)
	}
}

func SetLevelName(value string) {
	switch value {
	case "panic":
		_logger.SetLevel(logrus.PanicLevel)
		_level = LEVEL_PANIC
	case "error":
		_logger.SetLevel(logrus.ErrorLevel)
		_level = LEVEL_ERROR
	case "warn":
		_logger.SetLevel(logrus.WarnLevel)
		_level = LEVEL_WARN
	case "info":
		_logger.SetLevel(logrus.InfoLevel)
		_level = LEVEL_INFO
	case "debug":
		_logger.SetLevel(logrus.DebugLevel)
		_level = LEVEL_DEBUG
	case "trace":
		_logger.SetLevel(logrus.TraceLevel)
		_level = LEVEL_TRACE
	default:
		_logger.SetLevel(logrus.WarnLevel)
		_level = LEVEL_WARN
	}
}

func GetLevel() int {
	return _level
}

func SetOutput(value int) {
	_output = value
}

func Close() {
	if nil != _channel {
		close(_channel)
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	l o g g i n g
//----------------------------------------------------------------------------------------------------------------------

func Panic(args ...interface{}) {
	initialize()

	logParams := new(LogParams)
	logParams.Level = LEVEL_PANIC
	logParams.Args = args

	if nil != _channel {
		_channel <- logParams
	}
}

func Error(args ...interface{}) {
	initialize()

	logParams := new(LogParams)
	logParams.Level = LEVEL_ERROR
	logParams.Args = args
	if nil != _channel {
		_channel <- logParams
	}
}

func Warn(args ...interface{}) {
	initialize()

	logParams := new(LogParams)
	logParams.Level = LEVEL_WARN
	logParams.Args = args
	if nil != _channel {
		_channel <- logParams
	}
}

func Info(args ...interface{}) {
	initialize()

	logParams := new(LogParams)
	logParams.Level = LEVEL_INFO
	logParams.Args = args
	if nil != _channel {
		_channel <- logParams
	}
}

func Debug(args ...interface{}) {
	initialize()

	logParams := new(LogParams)
	logParams.Level = LEVEL_DEBUG
	logParams.Args = args
	if nil != _channel {
		_channel <- logParams
	}
}

func Trace(args ...interface{}) {
	initialize()

	logParams := new(LogParams)
	logParams.Level = LEVEL_TRACE
	logParams.Args = args
	if nil != _channel {
		_channel <- logParams
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func initialize() {
	if _initialized {
		return
	}
	_initialized = true

	// set logging file
	var workspace string = lygo_paths.GetWorkspacePath()
	_root = lygo_paths.Concat(workspace, ROOT) + string(os.PathSeparator)
	_fileName = lygo_paths.Concat(_root, "logging.log")

	// init the buffered channel
	_channel = make(chan *LogParams, 10)
	go receive(_channel)
}

func receive(ch <-chan *LogParams) {
	// loop until channel is open
	for arg := range ch {
		doLog(arg.Level, arg.Args...)
	}
}

func doLog(level int, args ...interface{}) {
	// ensure initialization
	initialize()

	// init write on file
	if _output == OUTPUT_FILE {
		if b, _ := lygo_paths.Exists(_root); !b {
			lygo_paths.Mkdir(_root)
		}

		file, err := os.OpenFile(_fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm) // 0644
		if err != nil {
			_logger.SetOutput(os.Stderr)
			_logger.Fatal(err)
		} else {
			// defer to close the file
			defer file.Close()
			_logger.SetOutput(file)
		}
	} else {
		_logger.SetOutput(os.Stdout)
	}

	// log context
	fields, message := parseContext(args...)
	// fields := logrus.Fields{}

	// do log
	switch level {
	case LEVEL_PANIC:
		_logger.WithFields(fields).Panic(message)
	case LEVEL_ERROR:
		_logger.WithFields(fields).Error(message)
	case LEVEL_WARN:
		_logger.WithFields(fields).Warn(message)
	case LEVEL_INFO:
		_logger.WithFields(fields).Info(message)
		// _logger.Info(message)
	case LEVEL_DEBUG:
		_logger.WithFields(fields).Debug(message)
	case LEVEL_TRACE:
		_logger.WithFields(fields).Trace(message)
	default:
		_logger.WithFields(fields).Info(message)

	}

}

func parseContext(args ...interface{}) (fields logrus.Fields, message string) {
	switch len(args) {
	case 1:
		message = fmt.Sprintf("%v", args[0])
		fields = logrus.Fields{}
	case 2:
		message = fmt.Sprintf("%v", args[1])
		context, ok := args[0].(LogContext)
		if ok {
			fields = logrus.Fields{
				"caller": context.Caller,
				"data":   context.Data,
			}
		} else {
			fields = logrus.Fields{}
		}
	default:
		message = fmt.Sprintf("%v", args[0])
		fields = logrus.Fields{}
	}

	return fields, message
}
