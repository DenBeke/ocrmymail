package ocrmymail

import (
	"fmt"
	"os/exec"
)

const cormypdfCommand = "ocrmypdf"

// OCRFile will use OCRmyPDF to perform OCR on the given input file
// and save the OCRed version to the give output file.
//
// This command expected OCRmyPDF to be installed.
// https://github.com/jbarlow83/OCRmyPDF
func OCRFile(input string, output string) error {

	cmd := exec.Command(cormypdfCommand, input, output)
	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("something went wrong while running ocrmypdf command: %w", err)
	}

	return nil

}

// IsOCRMyPDFInstalled checks wheter OCRmyPDF is installed.
func IsOCRMyPDFInstalled() bool {
	_, err := exec.LookPath(cormypdfCommand)
	if err == nil {
		return true
	}
	return false
}
