package contract

type PutMetadataRequest struct {
	DefaultDatasetID string `json:"default_dataset_id"`
}

type GetMetadataResponse struct {
	PopulationType   string `json:"population-type"`
	DefaultDatasetID string `json:"default_dataset_id"`
}
