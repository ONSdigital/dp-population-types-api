package contract

type PutMetadataRequest struct {
	PopulationType   string `schema:"population-type"`
	DefaultDatasetID string `json:"default_dataset_id"`
	Edition          string `json:"edition"`
	Version          int    `json:"version"`
}

type GetMetadataResponse struct {
	PopulationType   string `json:"population_type"`
	DefaultDatasetID string `json:"default_dataset_id"`
	Edition          string `json:"edition"`
	Version          int    `json:"version"`
}
