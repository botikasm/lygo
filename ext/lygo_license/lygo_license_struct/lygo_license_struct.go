package lygo_license_struct

import (
	"encoding/json"
	"github.com/botikasm/lygo/base/lygo_rnd"
	"github.com/botikasm/lygo/base/lygo_strings"
	"math"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

type License struct {
	Uid          string            `json:"uid"`
	CreationTime time.Time         `json:"creation_time"`
	DurationDays int64             `json:"duration_days"`
	Name         string            `json:"name"`
	Email        string            `json:"email"`
	Lang         string            `json:"lang"`
	Enabled      bool              `json:"enabled"`
	Params       map[string]string `json:"params"`
}

//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

func NewLicense(uid string) *License {
	instance := new(License)

	instance.CreationTime = time.Now() // seconds
	if len(uid) > 0 {
		instance.Uid = uid
	} else {
		instance.Uid = lygo_rnd.Uuid()
	}
	instance.Enabled = true
	instance.DurationDays = 1
	instance.Lang = "en"
	instance.Name = "Anonymous"
	instance.Email = ""

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	License
//----------------------------------------------------------------------------------------------------------------------

func (instance *License) Parse(text string) error {
	return json.Unmarshal([]byte(text), &instance)
}

func (instance *License) ToString() string {
	b, err := json.Marshal(&instance)
	if nil == err {
		return string(b)
	}
	return ""
}

func (instance *License) GetDataAsString() string {
	return lygo_strings.Format("\tId: %s\n\tOwner: %s\n\tCreation Date: %s\n\tExpire Date: %s\n\tExpired from days: %s",
		instance.Uid,
		instance.Name,
		instance.CreationTime,
		instance.GetExpireDate(),
		instance.RemainingDays()*-1,
	)
}

func (instance *License) IsValid() bool {
	if instance.Enabled && len(instance.Uid) > 0 {
		created := instance.CreationTime
		duration := instance.DurationDays
		now := time.Now()
		days := int64(now.Sub(created).Hours() / 24)

		return days <= duration
	}
	return false
}

func (instance *License) RemainingDays() int64 {
	if instance.Enabled && len(instance.Uid) > 0 {
		created := instance.CreationTime
		duration := instance.DurationDays
		now := time.Now()
		days := int64(now.Sub(created).Hours() / 24)

		return duration - days
	}
	return 0
}

func (instance *License) GetExpireDate() time.Time {
	created := instance.CreationTime
	duration := instance.DurationDays
	durationHours := duration * 24
	return created.Add(time.Duration(durationHours) * time.Hour)
}

func (instance *License) SetExpireDate(date time.Time) {
	created := instance.CreationTime
	days := int64(math.Round(date.Sub(created).Hours() / 24))
	instance.DurationDays = days
}

func (instance *License) ParseExpireDate(layout string, value string) error {
	date, err := time.Parse(layout, value)
	if nil == err {
		instance.SetExpireDate(date)
	}
	return err
}

func (instance *License) Add(days int64) {
	instance.Enabled = true
	instance.DurationDays = instance.DurationDays + days
}
