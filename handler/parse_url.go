package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type UrlQueryParam map[string]struct {
	is_url_param  bool   //is it URL param or a query param
	variable_type string //is variable is type string or int
}

func ParseRequest(r *http.Request, variables UrlQueryParam, req interface{}) error {
	logData := log.Data{}
	var m = make(map[string]interface{})
	for k, v := range variables {
		if v.variable_type == "int" {
			p := r.URL.Query().Get(k)
			if p != "" {
				logData[k] = p
				val, err := strconv.Atoi(p)
				if err != nil || val < 0 {
					err = errors.New("invalid query parameter")
					log.Error(r.Context(), "invalid query parameter: "+k, err, logData)
					return err
				}
				m[k] = val
			}
		} else {
			var val string
			if v.is_url_param {
				val = chi.URLParam(r, k)
				if val == "" {
					err := errors.New("invalid url parameter")
					log.Error(r.Context(), "invalid url parameter: "+k, err, logData)
					return err
				}

			} else { //query param can be blank
				val = r.URL.Query().Get(k)
			}
			m[k] = val
		}
	}
	c, _ := json.Marshal(m)
	json.Unmarshal(c, req)

	return nil
}
