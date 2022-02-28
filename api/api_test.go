//nolint:typecheck
package api_test

import (
	"github.com/gorilla/mux"
	"net/http/httptest"
)

//func TestSetup(t *testing.T) {
//	Convey("Given an API instance", t, func() {
//		r := mux.NewRouter()
//		ctx := context.Background()
//		api := api.Setup(ctx, r)
//
//		// TODO: remove hello world example handler route test case
//		Convey("When created the following routes should have been added", func() {
//			// Replace the check below with any newly added api endpoints
//			So(hasRoute(api.Router, "/hello", "GET"), ShouldBeTrue)
//		})
//	})
//}

func hasRoute(r *mux.Router, path, method string) bool {
	req := httptest.NewRequest(method, path, nil)
	match := &mux.RouteMatch{}
	return r.Match(req, match)
}
