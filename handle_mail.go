package ocrmymail

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"github.com/DusanKasan/parsemail"
	"github.com/gopistolet/smtp/smtp"
	"github.com/gosimple/slug"
	log "github.com/sirupsen/logrus"
)

// Handle implements the GoPistolet SMTP Handler interface
func (o *OCRMyMail) Handle(s *smtp.State) {

	log.Println("received an incoming mail...")
	// log.Printf("%+v", s)

	reader := bytes.NewReader(s.Data)

	email, err := parsemail.Parse(reader)
	if err != nil {
		log.Errorf("couldn't parse email body: %v", err)
	}

	log.Printf("%+v", email)

	attachmentsToMailOut := []string{}

	for _, attachment := range email.Attachments {
		if attachment.ContentType == "image/pdf" {

			documentName := getSafeDocumentName(attachment.Filename)
			attachmentFileName := tmpDir + documentName + ".pdf"

			err := downloadAttachment(attachmentFileName, attachment.Data)
			if err != nil {
				log.Errorf("couldn't download attachment: %v", err)
				return
			}
			log.Println("Saved attachment on disk")

			// We save the OCR version in the same file
			// that way on OCR error we send out at least the original file
			err = OCRFile(attachmentFileName, attachmentFileName)
			if err != nil {
				log.Errorln(err)
			}

			log.Println("Saved OCR version on disk")

			attachmentsToMailOut = append(attachmentsToMailOut, attachmentFileName)

		}
	}

	// Send the mail
	// TODO handle multiple recipients
	err = o.SendMail(s.To[0].String(), email.Subject, email.TextBody, attachmentsToMailOut)
	if err != nil {
		log.Errorf("couldn't mailout the OCRed content: %v", err)
	}

	fmt.Println("Relayed email to original recipient.")

}

func getSafeDocumentName(filename string) string {
	documentName := strings.TrimSuffix(filename, ".pdf")
	documentName = slug.Make(documentName)
	return documentName
}

func downloadAttachment(filename string, data io.Reader) error {

	attachmentBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return fmt.Errorf("coudldn't read attachment: %v", err)
	}

	// write the whole body at once
	err = ioutil.WriteFile(filename, attachmentBytes, 0644)
	if err != nil {
		return fmt.Errorf("couldn't write attachment to disk: %v", err)
	}

	return nil

}
