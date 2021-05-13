package main

import (
	log "github.com/sirupsen/logrus"

	ocrmymail "github.com/DenBeke/ocrmymail"
)

func main() {

	log.SetLevel(log.DebugLevel)

	config := ocrmymail.BuildConfigFromEnv()
	OCRMyMail, err := ocrmymail.New(config)
	if err != nil {
		log.Fatalf("couldn't create OCRMyMail instance: %v", err)
	}

	OCRMyMail.Serve()

}
