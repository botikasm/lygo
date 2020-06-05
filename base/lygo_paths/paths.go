package lygo_paths

import (
	"github.com/botikasm/lygo"
	"github.com/botikasm/lygo/base/lygo_regex"
	"github.com/botikasm/lygo/base/lygo_rnd"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	f i e l d s
//----------------------------------------------------------------------------------------------------------------------

var _workspace string = lygo.DEF_WORKSPACE
var _temp_root string = lygo.DEF_TEMP
var _pathSeparator = string(os.PathSeparator)

//----------------------------------------------------------------------------------------------------------------------
//	i n i t
//----------------------------------------------------------------------------------------------------------------------

func init() {
	_workspace = Absolute(lygo.DEF_WORKSPACE)
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func GetWorkspacePath() string {
	return _workspace
}

func SetWorkspacePath(value string) {
	_workspace = Absolute(value)
}

func SetWorkspaceParent(value string) {
	_workspace = filepath.Join(Absolute(value), lygo.DEF_WORKSPACE)
}

func WorkspacePath(partial string) string {
	if filepath.IsAbs(partial) {
		return partial
	}
	return filepath.Join(GetWorkspacePath(), partial)
}

func GetTempRoot() string {
	return Absolute(_temp_root)
}

func SetTempRoot(path string) {
	_temp_root = Absolute(path)
}

func Concat(paths ...string) string {
	return filepath.Join(paths...)
}

func ConcatDir(paths ...string) string {
	return filepath.Join(paths...) + string(os.PathSeparator)
}

// Check if a path exists and returns a boolean value or an error if access is denied
// @param path Path to check
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func Absolute(path string) string {
	abs, err := filepath.Abs(path)
	if nil == err {
		return abs
	}
	return path
}

func Extension(path string) string {
	return filepath.Ext(path)
}

func ExtensionName(path string) string {
	return strings.Replace(Extension(path), ".", "", 1)
}

func FileName(path string, includeExt bool) string {
	if IsUrl(path) {
		uri, err := url.Parse(path)
		if nil != err {
			return ""
		}
		path := uri.Path
		if len(path) > 1 {
			return FileName(path, includeExt)
		}
		return ""
	} else {
		base := filepath.Base(path)
		if includeExt {
			return base
		} else {
			ext := filepath.Ext(base)
			return strings.Replace(base, ext, "", 1)
		}
	}
}

// Creates a directory and all subdirectories if does not exists
func Mkdir(path string) (err error) {
	// ensure we have a directory
	var abs string
	if filepath.IsAbs(path) {
		abs = path
	} else {
		abs, err = filepath.Abs(path)
	}

	if nil == err {
		var f bool
		if f, _ = IsFile(abs); f {
			abs = filepath.Dir(abs)
		}

		if !strings.HasSuffix(abs, _pathSeparator) {
			path = abs + string(os.PathSeparator)
		} else {
			path = abs
		}

		var b bool
		if b, err = Exists(path); !b && nil == err {
			err = os.MkdirAll(path, os.ModePerm)
		}

	}

	return err
}

func IsTemp(path string) bool {
	tokens := strings.Split(path, _pathSeparator)
	temp := FileName(_temp_root, false)
	for _, token := range tokens {
		if token == temp {
			return true
		}
	}
	return false
}

func IsDir(path string) (bool, error) {
	fi, err := os.Lstat(Absolute(path))
	return fi.Mode().IsDir(), err
}

func IsFile(path string) (bool, error) {
	fi, err := os.Lstat(Absolute(path))
	if nil == err {
		return fi.Mode().IsRegular(), err
	} else {
		// path or file does not exists
		// just check if has extension
		if len(filepath.Ext(path)) > 0 {
			return true, nil
		}
		if strings.HasSuffix(path, _pathSeparator) {
			// is a directory
			return false, err
		}
		return true, err
	}
}

func IsAbs(path string) bool {
	if IsUrl(path) {
		return true
	}
	return filepath.IsAbs(path)
}

func IsUrl(path string) bool {
	return strings.Index(path, "http") == 0
}

func IsSymLink(path string) (bool, error) {
	fi, err := os.Lstat(Absolute(path))
	return fi.Mode()&os.ModeSymlink != 0, err
}

func IsHiddenFile(path string) (bool, error) {
	return isHiddenFile(path)
}

func IsSameFile(path1, path2 string) (bool, error) {
	f1, err1 := os.Lstat(Absolute(path1))
	f2, err2 := os.Lstat(Absolute(path2))
	if nil != err1 {
		return false, err1
	}
	if nil != err2 {
		return false, err2
	}
	return sameFile(f1, f2), nil
}

func IsSameFileInfo(f1, f2 os.FileInfo) bool {
	return sameFile(f1, f2)
}

func ListAll(root string) ([]string, error) {
	var response []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		response = append(response, path)
		return nil
	})
	return response, err
}

func ListFiles(root string, filter string) ([]string, error) {
	var response []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			if len(filter) == 0 {
				response = append(response, path)
			} else {
				name := filepath.Base(path)
				if len(lygo_regex.WildcardMatch(name, filter)) > 0 {
					response = append(response, path)
				}
			}
		}
		return nil
	})
	return response, err
}

func TmpFileName(extension string) string {
	uuid := lygo_rnd.Uuid()
	if len(uuid) == 0 {
		uuid = "temp_file"
	}

	return uuid + ensureDot(extension)
}

func TmpFile(extension string) string {
	path := filepath.Join(_temp_root, TmpFileName(extension))
	return Absolute(path)
}

func ChangeFileName(fromPath, toFileName string) string {
	parent := filepath.Dir(fromPath)
	base := filepath.Base(fromPath)
	ext := filepath.Ext(base)
	if len(filepath.Ext(toFileName)) > 0 {
		return filepath.Join(parent, toFileName)
	}
	return filepath.Join(parent, toFileName+ext)
}

func ChangeFileNameExtension(fromPath, toFileExtension string) string {
	parent := filepath.Dir(fromPath)
	base := filepath.Base(fromPath)
	ext := filepath.Ext(base)
	name := strings.Replace(base, ext, "", 1)
	return filepath.Join(parent, name+ensureDot(toFileExtension))
}

func ChangeFileNameWithSuffix(fileName, suffix string) string {
	base := filepath.Base(fileName)
	ext := filepath.Ext(base)
	name := strings.Replace(base, ext, "", 1)
	return filepath.Join(filepath.Dir(fileName), name+suffix+ext)
}

func ChangeFileNameWithPrefix(fileName, prefix string) string {
	base := filepath.Base(fileName)
	return filepath.Join(filepath.Dir(fileName), prefix+base)
}

func CleanPath(p string) string {
	return filepath.Clean(p)
}

// CleanUrl is the URL version of path.Clean, it returns a canonical URL path
// for p, eliminating . and .. elements.
//
// The following rules are applied iteratively until no further processing can
// be done:
//	1. Replace multiple slashes with a single slash.
//	2. Eliminate each . path name element (the current directory).
//	3. Eliminate each inner .. path name element (the parent directory)
//	   along with the non-.. element that precedes it.
//	4. Eliminate .. elements that begin a rooted path:
//	   that is, replace "/.." by "/" at the beginning of a path.
//
// If the result of this process is an empty string, "/" is returned
func CleanUrl(p string) string {
	// Turn empty string into "/"
	if p == "" {
		return "/"
	}

	n := len(p)
	var buf []byte

	// Invariants:
	//      reading from path; r is index of next byte to process.
	//      writing to buf; w is index of next byte to write.

	// path must start with '/'
	r := 1
	w := 1

	if p[0] != '/' {
		r = 0
		buf = make([]byte, n+1)
		buf[0] = '/'
	}

	trailing := n > 2 && p[n-1] == '/'

	// A bit more clunky without a 'lazybuf' like the path package, but the loop
	// gets completely inlined (bufApp). So in contrast to the path package this
	// loop has no expensive function calls (except 1x make)

	for r < n {
		switch {
		case p[r] == '/':
			// empty path element, trailing slash is added after the end
			r++

		case p[r] == '.' && r+1 == n:
			trailing = true
			r++

		case p[r] == '.' && p[r+1] == '/':
			// . element
			r++

		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'):
			// .. element: remove to last /
			r += 2

			if w > 1 {
				// can backtrack
				w--

				if buf == nil {
					for w > 1 && p[w] != '/' {
						w--
					}
				} else {
					for w > 1 && buf[w] != '/' {
						w--
					}
				}
			}

		default:
			// real path element.
			// add slash if needed
			if w > 1 {
				bufApp(&buf, p, w, '/')
				w++
			}

			// copy element
			for r < n && p[r] != '/' {
				bufApp(&buf, p, w, p[r])
				w++
				r++
			}
		}
	}

	// re-append trailing slash
	if trailing && w > 1 {
		bufApp(&buf, p, w, '/')
		w++
	}

	if buf == nil {
		return p[:w]
	}
	return string(buf[:w])
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func ensureDot(extension string) string {
	if strings.Index(extension, ".") == -1 {
		extension = "." + extension
	}
	return extension
}

// internal helper to lazily create a buffer if necessary
func bufApp(buf *[]byte, s string, w int, c byte) {
	if *buf == nil {
		if s[w] == c {
			return
		}

		*buf = make([]byte, len(s))
		copy(*buf, s[:w])
	}
	(*buf)[w] = c
}
