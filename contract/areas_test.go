package contract_test

import (
	"testing"

	"github.com/ONSdigital/dp-population-types-api/contract"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetAreasRequestValid(t *testing.T) {
	var limit20 = 20
	var negativeLimit = -2
	Convey("Given a valid GetAreasRequest object", t, func() {
		req := contract.GetAreasRequest{
			QueryParams: contract.QueryParams{
				Limit:  &limit20,
				Offset: 0,
			},
			Category: "hello",
		}

		Convey("When Valid() is called", func() {
			err := req.Valid()
			So(err, ShouldBeNil)
		})
	})

	Convey("Given an invalid value for limit is given", t, func() {
		req := contract.GetAreasRequest{
			QueryParams: contract.QueryParams{
				Limit:  &negativeLimit,
				Offset: 0,
			},
			Category: "hello",
		}

		Convey("When Valid() is called", func() {
			err := req.Valid()
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given an invalid value for offset is given", t, func() {
		req := contract.GetAreasRequest{
			QueryParams: contract.QueryParams{
				Limit:  &limit20,
				Offset: -10,
			},
			Category: "hello",
		}

		Convey("When Valid() is called", func() {
			err := req.Valid()
			So(err, ShouldNotBeNil)
		})
	})

	Convey("Given an URL encoded value is given for Category", t, func() {
		req := contract.GetAreasRequest{
			QueryParams: contract.QueryParams{
				Limit:  &limit20,
				Offset: 0,
			},
			Category: "hello%20there",
		}

		Convey("When Valid() is called", func() {
			err := req.Valid()
			So(err, ShouldBeNil)

			Convey("Category should resolve to the decoded value", func() {
				So(req.Category, ShouldEqual, "hello there")
			})
		})
	})

	Convey("Given an incorecctly encoded URL value is given for Category", t, func() {
		req := contract.GetAreasRequest{
			QueryParams: contract.QueryParams{
				Limit:  &limit20,
				Offset: 0,
			},
			Category: "North%&@w20East",
		}

		Convey("When Valid() is called", func() {
			err := req.Valid()
			So(err, ShouldNotBeNil)
		})
	})
}
