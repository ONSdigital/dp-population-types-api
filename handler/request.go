package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/ONSdigital/log.go/v2/log"
	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

// parseRequest attemts to read unmarshal a request body and read URL
// params into a request object, returning an appropriate error on failure.
// 'req' must be a pointer to a struct.
// ParseRequest will also attempt to call the request's Valid()
// function if it has one and will throw an error if it fails
func parseRequest(r *http.Request, req interface{}) error {
	var b []byte
	if r.Body != nil && r.Body != http.NoBody {
		b, err := io.ReadAll(r.Body)
		if err != nil {
			return Error{
				err:     errors.Wrap(err, "failed to read request body"),
				message: "failed to read request body",
			}
		}

		// parse request body
		if len(b) > 0 {
			if err := json.Unmarshal(b, &req); err != nil {
				return Error{
					err:     fmt.Errorf("failed to unmarshal request body: %w", err),
					message: "badly formed request body",
					logData: log.Data{
						"body": string(b),
					},
				}
			}
		}
	}

	// parse URL query params
	if err := schema.NewDecoder().Decode(req, r.URL.Query()); err != nil {
		return Error{
			err:     errors.Wrap(err, "failed to decode query parameters"),
			message: "badly formed query parameters",
			logData: log.Data{
				"query_params": r.URL.Query(),
			},
		}
	}

	if v, ok := req.(validator); ok {
		if err := v.Valid(); err != nil {
			return Error{
				err: errors.Wrap(err, "invalid request"),
				logData: log.Data{
					"body":    string(b),
					"request": fmt.Sprintf("%+v", req),
				},
			}
		}
	}

	return nil
}
