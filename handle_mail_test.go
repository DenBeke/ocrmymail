package ocrmymail

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetSafeDocumentName(t *testing.T) {

	Convey("getSafeDocumentName", t, func() {

		data := []struct {
			input  string
			output string
		}{
			{
				"test",
				"test",
			},
			{
				"te/st",
				"te-st",
			},
			{
				"tést",
				"test",
			},
			{
				"tést!",
				"test",
			},
			{
				"te.st",
				"te-st",
			},
		}

		for _, test := range data {

			So(getSafeDocumentName(test.input), ShouldEqual, test.output)

		}

	})
}
