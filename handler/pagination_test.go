package handler

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestManualPagination(t *testing.T) {

	testInput := []string{
		"1", "2", "3", "4",
	}
	Convey("Given a set of parameters for pagination", t, func() {
		Convey("When non integer characters are given", func() {
			offset, limit := "wrong", "wrong"
			output, limitOutput, offsetOutput := manualPagination(limit, offset, testInput)

			So(len(output.Items), ShouldEqual, 4)
			So(limitOutput, ShouldEqual, 4)
			So(offsetOutput, ShouldEqual, 0)

		})

		Convey("When reasonable parameters are given for offset and limit", func() {
			offset, limit := "0", "20"
			output, limitOutput, offsetOutput := manualPagination(limit, offset, testInput)

			So(len(output.Items), ShouldEqual, 4)
			So(limitOutput, ShouldEqual, 4)
			So(offsetOutput, ShouldEqual, 0)

		})

		Convey("When unreasonable parameters are given for offset and limit", func() {
			offset, limit := "100", "100"
			output, limitOutput, offsetOutput := manualPagination(limit, offset, testInput)

			So(len(output.Items), ShouldEqual, 4)
			So(limitOutput, ShouldEqual, 4)
			So(offsetOutput, ShouldEqual, 0)

		})
	})
}
