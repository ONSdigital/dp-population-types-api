package contract

type PaginationResponse struct {
	Limit      *int `json:"limit,omitempty"`
	Offset     int  `json:"offset"`
	Count      int  `json:"count"`
	TotalCount int  `json:"total_count"`
}
