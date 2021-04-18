package ocrmymail

import (
	"github.com/pkg/errors"
	"gopkg.in/mail.v2"
)

// SendMail sends a mail with the given recipient, subject body and attachments.
// Attachments should contain a list of filenames.
func (m *OCRMyMail) SendMail(to string, subject string, body string, attachments []string) error {

	// Construct mail message
	msg := mail.NewMessage()
	msg.SetHeader("From", m.config.SMTP.FromEmail)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	for _, filename := range attachments {
		msg.Attach(filename)
	}

	d := mail.NewDialer(
		m.config.RemoteSMTP.Host,
		m.config.RemoteSMTP.Port,
		m.config.RemoteSMTP.User,
		m.config.RemoteSMTP.Password,
	)

	if m.config.RemoteSMTP.DisableTLS {
		d.StartTLSPolicy = mail.NoStartTLS
	}

	// Send the actual mail
	if err := d.DialAndSend(msg); err != nil {
		return errors.Wrap(err, "couldn't send the email")
	}

	return nil
}
