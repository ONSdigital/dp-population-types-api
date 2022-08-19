package contract

import "github.com/pkg/errors"

const (
	defaultLimit = 20
)

type QueryParams struct {
	Limit  int `schema:"limit"`
	Offset int `schema:"offset"`
}

func (q *QueryParams) Valid() error {
	if q.Limit < 0 {
		return errors.New("'limit' parameter must not be less than 0")
	}

	if q.Limit == 0 {
		q.Limit = defaultLimit
	}

	if q.Offset < 0 {
		return errors.New("'offset' parameter must not be less than 0")
	}

	return nil
}
