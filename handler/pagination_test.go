package handler

import (
	"testing"

	"github.com/ONSdigital/dp-population-types-api/contract"
	. "github.com/smartystreets/goconvey/convey"
)

func TestManualPagination(t *testing.T) {

	testInput := []contract.PopulationType{
		{
			Name: "1",
		},
		{
			Name: "2",
		},
		{
			Name: "3",
		},
		{
			Name: "4",
		},
	}
	Convey("Given a set of parameters for pagination", t, func() {

		Convey("When reasonable parameters are given for offset and limit", func() {
			offset, limit := 0, 20
			r := contract.GetPopulationTypesResponse{
				PaginationResponse: contract.PaginationResponse{
					Limit:  limit,
					Offset: offset,
				},
				Items: testInput,
			}
			r.Paginate()

			So(len(r.Items), ShouldEqual, 4)
		})

		Convey("When unreasonable parameters are given for offset and limit", func() {
			offset, limit := 100, 100
			r := contract.GetPopulationTypesResponse{
				PaginationResponse: contract.PaginationResponse{
					Limit:  limit,
					Offset: offset,
				},
				Items: testInput,
			}
			r.Paginate()

			So(len(r.Items), ShouldEqual, 4)

		})
	})
}
