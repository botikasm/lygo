package lygo_license_client

import (
	"errors"
	"github.com/botikasm/lygo/base/lygo_event_ticker"
	"github.com/botikasm/lygo/base/lygo_regex"
	"github.com/botikasm/lygo/base/lygo_strings"
	"github.com/botikasm/lygo/ext/lygo_license/lygo_license_config"
	"github.com/botikasm/lygo/ext/lygo_license/lygo_license_struct"
	"net/smtp"
	"strings"
	"time"
)

//----------------------------------------------------------------------------------------------------------------------
//	t y p e s
//----------------------------------------------------------------------------------------------------------------------

/**
 * Timed license monitor.
 * This object check license validity every X hours.
 */
type LicenseTicker struct {

	// ticker interval
	Interval time.Duration
	// optional path for license (usually is defined in configuration)
	LicensePath string

	// [run only if expired] override default action (send email if expired)
	ExpiredLicenseHook LicenseTickerCallback
	// [run always] notify action on license
	RequestLicenseHook LicenseTickerCallback

	Email *LicenseTickerEmail
	// count number of time the license expire email was sent
	CountWarnings int

	//-- private --//
	ticker  *lygo_event_ticker.EventTicker
	client  *LicenseClient
	stopped bool
}

type LicenseTickerContext struct {
	Error   error
	License *lygo_license_struct.License
	Ticker  *LicenseTicker
}
func (instance *LicenseTickerContext) HasError404() bool {
	if nil != instance.Error {
		return strings.Index(instance.Error.Error(), "404") > -1
	}
	return false
}

type LicenseTickerCallback func(*LicenseTickerContext)

type LicenseTickerEmail struct {
	Enabled bool

	SmtpHost     string // smtp.gmail.com
	SmtpPort     int    // 587
	SmtpUser     string
	SmtpPassword string

	From    string
	Subject string
	Message string

	Target []string

	Errors []error

	//-- private --//
	sent time.Time
}

func (instance *LicenseTickerEmail) CanSend() bool {
	last := instance.sent
	now := time.Now()
	diff := now.Sub(last)
	hours := diff.Hours()
	return hours > 12
}


//----------------------------------------------------------------------------------------------------------------------
//	c o n s t r u c t o r
//----------------------------------------------------------------------------------------------------------------------

/*
 * Creates new License Ticker.
 * @param config
 * @param callback Pass is want override default (send email to license owner)
 */
func NewLicenseTicker(config *lygo_license_config.LicenseConfig) *LicenseTicker {
	instance := new(LicenseTicker)
	instance.client = NewLicenseClient(config)

	instance.Email = new(LicenseTickerEmail)
	instance.Email.Enabled = true

	instance.stopped = false
	instance.Interval = 1 * time.Hour
	instance.LicensePath = ""

	return instance
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *LicenseTicker) Start() {
	if nil == instance.ticker {
		instance.stopped = false

		// initial check
		instance.doCheck()

		if !instance.stopped {
			// enable ticker
			instance.ticker = lygo_event_ticker.NewEventTicker(instance.Interval, instance.onTick)
		}
	}
}

func (instance *LicenseTicker) Stop() {
	instance.stopped = true
	if nil != instance.ticker {
		instance.ticker.Stop()
		instance.ticker = nil
	}
}

func (instance *LicenseTicker) Join() {
	if nil != instance.ticker && !instance.stopped {
		instance.ticker.Join()
	}
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *LicenseTicker) onTick(ticker *lygo_event_ticker.EventTicker) {
	instance.doCheck()
}

func (instance *LicenseTicker) doCheck() {
	context := new(LicenseTickerContext)
	context.Ticker = instance

	// request the license
	if nil != instance.client {
		license, err := instance.client.RequestLicense(instance.LicensePath)
		context.License = license
		context.Error = err
		if nil != err {

		}
	} else {
		context.Error = errors.New("system: License client is null")
	}

	if nil != context.License {
		// run only if license is expired
		if !context.License.IsValid() {
			instance.CountWarnings++
			if nil == instance.ExpiredLicenseHook {
				go instance.doActions(context)
			} else {
				instance.ExpiredLicenseHook(context)
			}
		} else {
			instance.CountWarnings = 0
		}

		// call hook only if context has no errors
		if nil != instance.RequestLicenseHook {
			instance.RequestLicenseHook(context)
		}
	}
}

func (instance *LicenseTicker) doActions(context *LicenseTickerContext) {
	if instance.Email.Enabled && nil == context.Error && instance.Email.CanSend() {
		instance.Email.sent = time.Now() //avoid spam
		instance.Email.Errors = make([]error, 0)

		if len(instance.Email.Subject) > 0 && len(instance.Email.Message) > 0 && len(instance.Email.SmtpHost) > 0 {
			// send email to license owner
			name := context.License.Name
			email := context.License.Email
			data := context.License.GetDataAsString()

			if len(email) > 0 {
				instance.Email.Target = append(instance.Email.Target, email)

				subject := lygo_strings.Format(instance.Email.Subject, name)
				body := lygo_strings.Format(instance.Email.Message, name, data)

				user := instance.Email.SmtpUser
				psw := instance.Email.SmtpPassword
				host := instance.Email.SmtpHost
				port := instance.Email.SmtpPort
				address := lygo_strings.Format("%s:%s", host, port)
				from := instance.Email.From
				if len(from) == 0 {
					from = user
				}
				auth := smtp.PlainAuth("", user, psw, host)

				for _, email := range instance.Email.Target {
					if len(lygo_regex.Emails(email)) > 0 {
						msg := "From: " + from + "\n" +
							"To: " + email + "\n" +
							"Subject: " + subject + "\n\n" +
							body

						err := smtp.SendMail(address, auth, from, []string{email}, []byte(msg))
						if nil != err {
							instance.Email.Errors = append(instance.Email.Errors, err)
							// stop loop if error occurred
							break
						}
					}
				}
			}
		}
	}
}
