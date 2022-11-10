package handler

import (
	"testing"

	"github.com/ONSdigital/dp-population-types-api/contract"
	. "github.com/smartystreets/goconvey/convey"
)

func TestManualPagination(t *testing.T) {

	testInput := []string{
		"1", "2", "3", "4",
	}
	Convey("Given a set of parameters for pagination", t, func() {

		Convey("When reasonable parameters are given for offset and limit", func() {
			offset, limit := 0, 20
			r := contract.GetPopulationTypesRequest{
				QueryParams: contract.QueryParams{
					Limit:  limit,
					Offset: offset,
				},
			}
			output := r.Paginate(testInput)

			So(len(output.Items), ShouldEqual, 4)

		})

		Convey("When unreasonable parameters are given for offset and limit", func() {
			offset, limit := 100, 100
			r := contract.GetPopulationTypesRequest{
				QueryParams: contract.QueryParams{
					Limit:  limit,
					Offset: offset,
				},
			}

			output := r.Paginate(testInput)

			So(len(output.Items), ShouldEqual, 4)

		})
	})
}
