package lygo_license

import (
	"fmt"
	"github.com/botikasm/lygo/base/lygo_io"
	"github.com/botikasm/lygo/ext/lygo_license/lygo_license_client"
	"github.com/botikasm/lygo/ext/lygo_license/lygo_license_config"
	"github.com/botikasm/lygo/ext/lygo_license/lygo_license_struct"
	"testing"
	"time"
)

func TestSimple(t *testing.T) {

	text, err := lygo_io.ReadTextFromFile("./lygo_license_config.json")
	if nil == err {
		// creates client
		config := new(lygo_license_config.LicenseConfig)
		config.Parse(text)

		if len(config.Host) > 0 {
			client := lygo_license_client.NewLicenseClient(config)
			if nil != client {
				filename := client.Config.GetRequestFileName()
				if len(filename) > 0 {
					// have a direct file
					fmt.Println(filename)
				}
				license, err := client.RequestLicense("")
				if nil != err {
					t.Error(err)
				} else {
					fmt.Println(license)
					fmt.Println("CREATED: ", license.CreationTime)
					fmt.Println("DAYS: ", license.DurationDays)
					fmt.Println("EXPIRED: ", !license.IsValid())
					fmt.Println("EXPIRE DATE: ", license.GetExpireDate())
					fmt.Println("REMAINING: ", license.RemainingDays())
				}
			}
		} else {
			t.Error("Mismatch Configuration: Missing 'Host'")
			t.Fail()
		}
	} else {
		t.Error(err)
		t.Fail()
	}
}

func TestLicenseStruct(t *testing.T) {

	text, err := lygo_io.ReadTextFromFile("./lygo_license_struct.json")
	if nil == err {
		license := new(lygo_license_struct.License)
		license.Parse(text)

		// lygo_io.WriteTextToFile(license.ToString(), "./lygo_license_struct.json")

		if license.IsValid() {
			fmt.Println("valid license")
		} else {
			remaining := license.RemainingDays()
			fmt.Println("expired license... adding days", remaining*-1)
			license.Add(remaining * -1)
			fmt.Println("valid license", license.IsValid())

			// set expire data
			license.ParseExpireDate("2006-01-02T15:04:05.000Z", "2020-01-31T15:04:05.000Z")
			fmt.Println("NEW EXPIRE DATE: ", license.GetExpireDate())
			fmt.Println("EXPIRED: ", !license.IsValid())
			fmt.Println("EXPIRED DAYS: ", license.RemainingDays())
		}
	} else {
		t.Error(err)
	}

}

func TestLicenseTicker(t *testing.T) {

	text, err := lygo_io.ReadTextFromFile("./lygo_license_config.json")
	if nil == err {
		config := new(lygo_license_config.LicenseConfig)
		config.Parse(text)
		ticker := lygo_license_client.NewLicenseTicker(config)
		ticker.RequestLicenseHook = onLicense
		ticker.Email.Enabled = true
		ticker.Email.From = "Botika<info@botika.ai>"
		ticker.Email.Subject = "%s, License Expired"
		ticker.Email.Message = "Hi %s, \nyour license is expired \n%s"
		ticker.Email.SmtpHost = "ssl0.ovh.net"
		ticker.Email.SmtpPort = 587
		ticker.Email.SmtpUser = "support@botika.it"
		ticker.Email.SmtpPassword = "Fa1%"
		ticker.Email.Target = []string{}

		ticker.Start()

		// lock and wait manual stop
		ticker.Join()
	} else {
		t.Error(err)
	}

	// wait 4 seconds to allow email send
	fmt.Println("EXITING....")
	time.Sleep(15 * time.Second)
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func onLicense(ctx *lygo_license_client.LicenseTickerContext) {
	if nil!=ctx.License{
		license := ctx.License
		fmt.Println("NEW EXPIRE DATE: ", license.GetExpireDate())
		fmt.Println("EXPIRED: ", !license.IsValid())
		fmt.Println("EXPIRED DAYS: ", license.RemainingDays())

		ctx.Ticker.Stop()
	}
}
