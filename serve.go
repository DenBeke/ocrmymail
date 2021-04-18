package ocrmymail

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/gopistolet/smtp/mta"
	log "github.com/sirupsen/logrus"
)

const tmpDir = "./tmp/"

func (o *OCRMyMail) Serve() {

	// Validate config
	err := o.config.Validate()
	if err != nil {
		log.Fatalf("Config file is not valid: %v", err)
	}

	log.WithField("config", fmt.Sprintf("%+v", o.config)).Println("Starting PDF OCR SMTP Gateway ✉️")

	if !IsOCRMyPDFInstalled() {
		log.Fatalln("OCRmyPDF is not installed. Please install the command: https://github.com/jbarlow83/OCRmyPDF")
	}

	// Configure and start GoPistolet SMTP server
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, os.Interrupt, syscall.SIGTERM)

	// Default config
	smtpConfig := mta.Config{
		Hostname: "localhost",
		Port:     25,
	}

	// create new MTA with SMTP config and OCRMyMail as the email handler
	mta := mta.NewDefault(smtpConfig, o)
	go func() {
		<-sigc
		mta.Stop()
	}()
	err = mta.ListenAndServe()
	if err != nil {
		log.Fatalln(err)
	}
}
