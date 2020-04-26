package lygo_json

import (
	"encoding/json"
	"github.com/botikasm/lygo/base/lygo_io"
)

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func Bytes(entity interface{}) []byte {
	b, err := json.Marshal(&entity)
	if nil == err {
		return b
	}
	return []byte{}
}

func Stringify(entity interface{}) string {
	return string(Bytes(entity))
}

func Read(input interface{}, entity interface{}) (err error) {
	if s, b := input.(string); b {
		err = json.Unmarshal([]byte(s), &entity)
	} else if s, b := input.([]byte); b {
		err = json.Unmarshal(s, &entity)
	}
	if nil != err {
		return err
	}
	return nil
}

func ReadFromFile(fileName string, entity interface{}) error {
	b, err := lygo_io.ReadBytesFromFile(fileName)
	if nil != err {
		return err
	}
	err = json.Unmarshal(b, &entity)
	if nil != err {
		return err
	}
	return nil
}

func ReadMapFromFile(fileName string) (map[string]interface{}, error) {
	b, err := lygo_io.ReadBytesFromFile(fileName)
	if nil != err {
		return nil, err
	}
	var response map[string]interface{}
	err = json.Unmarshal(b, &response)
	if nil != err {
		return nil, err
	}
	return response, nil
}

func ReadArrayFromFile(fileName string) ([]map[string]interface{}, error) {
	b, err := lygo_io.ReadBytesFromFile(fileName)
	if nil != err {
		return nil, err
	}
	var response []map[string]interface{}
	err = json.Unmarshal(b, &response)
	if nil != err {
		return nil, err
	}
	return response, nil
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------
