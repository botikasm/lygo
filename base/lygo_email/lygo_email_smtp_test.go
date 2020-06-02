package lygo_email

import (
	"net/mail"
	"net/smtp"
	"testing"
)

func TestSimple(t *testing.T) {

	username := "xxx"
	password := "xxx"

	// compose the message
	m := NewMessage("Hi", "this is the body")
	m.From = mail.Address{Name: "From", Address: "from@example.com"}
	m.AddTo(mail.Address{Name: "To", Address: "xxxx@example.com"})
	m.AddCc(mail.Address{Name: "someCcName", Address: "xxxx@example.com"})
	m.AddBcc(mail.Address{Name: "someBccName", Address: "xxxx@example.com"})

	// add attachments
	if err := m.AddAttachment("readme.md"); err != nil {
		t.Error(err)
		t.FailNow()
	}

	// add headers
	m.AddHeader("X-CUSTOMER-id", "1234567789")

	// send it
	auth := smtp.PlainAuth("", username, password, "ssl0.ovh.net")
	if err := Send("ssl0.ovh.net:587", auth, m); err != nil {
		t.Error(err)
		t.FailNow()
	}

}