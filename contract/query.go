package contract

type QueryParams struct {
	Limit  int `schema:"limit"`
	Offset int `schema:"offset"`
}
