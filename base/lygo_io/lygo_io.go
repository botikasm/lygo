package lygo_io

import (
	"bufio"
	"fmt"
	"github.com/botikasm/lygo/base/lygo_paths"
	"io"
	"io/ioutil"
	"os"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func Remove(filename string) error {
	return os.Remove(filename)
}

func RemoveAll(path string) error {
	return os.RemoveAll(path)
}

func CopyFile(src, dst string) (int64, error) {
	lygo_paths.Mkdir(dst)

	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func AppendTextToFile(text, file string) (bytes int, err error) {
	var f *os.File
	if b, _ := lygo_paths.Exists(file); b {
		f, err = os.OpenFile(file,
			os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	} else {
		f, err = os.Create(file)
	}

	if nil == err {
		defer f.Close()
		w := bufio.NewWriter(f)
		bytes, err = w.WriteString(text)
		w.Flush()
	}
	return bytes, err
}

func WriteTextToFile(text, file string) (bytes int, err error) {
	var f *os.File
	f, err = os.Create(file)

	if nil == err {
		defer f.Close()
		w := bufio.NewWriter(f)
		bytes, err = w.WriteString(text)
		w.Flush()
	}
	return bytes, err
}

func WriteBytesToFile(data []byte, file string) (bytes int, err error) {
	var f *os.File
	f, err = os.Create(file)
	if nil == err {
		defer f.Close()
		w := bufio.NewWriter(f)
		bytes, err = w.Write(data)
		w.Flush()
	}
	return bytes, err
}

func ReadBytesFromFile(fileName string) ([]byte, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	b, err := ioutil.ReadAll(file)
	return b, err
}

func ReadTextFromFile(fileName string) (string, error) {
	b, err := ReadBytesFromFile(fileName)
	return string(b), err
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
