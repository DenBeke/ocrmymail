package ocrmymail

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/DusanKasan/parsemail"
	"github.com/google/uuid"
	"github.com/gopistolet/smtp/smtp"
	"github.com/gosimple/slug"
	log "github.com/sirupsen/logrus"
)

// Handle implements the GoPistolet SMTP Handler interface
func (o *OCRMyMail) Handle(s *smtp.State) {

	if o.config.HandleAsync {
		log.Debugln("handling incoming mail async")
		state := *s
		go o.handleMail(&state)
	} else {
		log.Debugln("handling incoming mail sync")
		o.handleMail(s)
	}

}

func (o *OCRMyMail) handleMail(s *smtp.State) {

	log.Println("parsing incoming mail...")
	//log.Debugf("%+v", s)

	reader := bytes.NewReader(s.Data)

	email, err := parsemail.Parse(reader)
	if err != nil {
		log.Errorf("couldn't parse email body: %v", err)
		return
	}

	log.Printf("%+v", email)

	attachmentsToMailOut := []Attachment{}

	for _, attachment := range email.Attachments {
		if attachment.ContentType == "image/pdf" {

			documentName := getSafeDocumentName(attachment.Filename)
			attachmentFileName := tmpDir + documentName + "_" + uuid.NewString() + ".pdf"

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

			attachmentsToMailOut = append(attachmentsToMailOut, Attachment{FilenameDisk: attachmentFileName, FilenameMail: documentName + ".pdf"})

		}
	}

	// Send the mail
	// TODO handle multiple recipients
	err = o.SendMail(s.To[0].String(), email.Subject, email.TextBody, attachmentsToMailOut)
	if err != nil {
		log.Errorf("couldn't mailout the OCRed content: %v", err)
	}

	// Delete attachments after they are mailed out
	for _, attachment := range attachmentsToMailOut {
		err := os.Remove(attachment.FilenameDisk)
		if err != nil {
			log.Errorf("couldn't delete attachment: %v", err)
		}
	}

	log.Println("Relayed email to original recipient.")

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
