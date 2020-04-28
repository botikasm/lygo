package lygo_db_sync

import (
	"encoding/json"
	"github.com/botikasm/lygo/base/lygo_conv"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/base/lygo_strings"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e
//----------------------------------------------------------------------------------------------------------------------

type DBSyncConfig struct {
	Address  string                `json:"address"`
	Database *DBSyncDatabaseConfig `json:"database"`
	Sync     []*DBSyncConfigSync   `json:"sync"`
}

type DBSyncDatabaseConfig struct {
	Endpoints      []string                  `json:"endpoints"`
	Authentication *DBSyncDatabaseConfigAuth `json:"authentication"`
}

type DBSyncDatabaseConfigAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type DBSyncConfigSync struct {
	LocalDBName  string `json:"local_dbname"`
	RemoteDBName string `json:"remote_dbname"`
}

//----------------------------------------------------------------------------------------------------------------------
//	DBSyncConfig
//----------------------------------------------------------------------------------------------------------------------

func (instance *DBSyncConfig) Load(path string) error {
	text, err := lygo_io.ReadTextFromFile(path)
	if nil != err {
		return err
	}
	return instance.Parse(text)
}

func (instance *DBSyncConfig) Parse(text string) error {
	return json.Unmarshal([]byte(text), &instance)
}

func (instance *DBSyncConfig) ToString() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}

func (instance *DBSyncConfig) Host() string {
	if len(instance.Address) > 0 {
		tokens := lygo_strings.Split(instance.Address, ":")
		if len(tokens) > 0 {
			return tokens[0]
		}
	}
	return ""
}

func (instance *DBSyncConfig) Port() int {
	if len(instance.Address) > 0 {
		tokens := lygo_strings.Split(instance.Address, ":")
		if len(tokens) == 2 {
			return lygo_conv.ToInt(tokens[1])
		} else {
			return lygo_conv.ToInt(tokens[0])
		}
	}
	return -1
}

//----------------------------------------------------------------------------------------------------------------------
//	DBSyncDatabaseConfig
//----------------------------------------------------------------------------------------------------------------------

func (instance *DBSyncDatabaseConfig) Parse(text string) error {
	return json.Unmarshal([]byte(text), &instance)
}

func (instance *DBSyncDatabaseConfig) ToString() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}

//----------------------------------------------------------------------------------------------------------------------
//	DBSyncDatabaseConfigAuth
//----------------------------------------------------------------------------------------------------------------------

func (instance *DBSyncDatabaseConfigAuth) Parse(text string) error {
	return json.Unmarshal([]byte(text), &instance)
}

func (instance *DBSyncDatabaseConfigAuth) ToString() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}
