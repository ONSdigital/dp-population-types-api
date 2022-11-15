package contract

import "github.com/pkg/errors"

var defaultLimit = 20

type QueryParams struct {
	Limit  *int `schema:"limit, omitempty"`
	Offset int  `schema:"offset"`
}

func (q *QueryParams) Valid() error {
	if q.Limit == nil {
		q.Limit = &defaultLimit
	}

	if *q.Limit < 0 {
		return errors.New("'limit' parameter must not be less than 0")
	}

	if q.Offset < 0 {
		return errors.New("'offset' parameter must not be less than 0")
	}

	return nil
}
