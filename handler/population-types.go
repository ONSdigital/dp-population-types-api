package handler

import (
	"net/http"
)

type PopulationTypes struct {
	responder responder
}

func NewPopulationTypes(responder responder) *PopulationTypes {
	return &PopulationTypes{
		responder: responder,
	}
}

type blah struct {
	Items []string `json:"items"`
}

func (h *PopulationTypes) Get(w http.ResponseWriter, req *http.Request) {
	body := blah{
		Items: make([]string, 0),
	}
	responder := h.responder
	ctx := req.Context()

	responder.JSON(ctx, w, 200, body)
}
