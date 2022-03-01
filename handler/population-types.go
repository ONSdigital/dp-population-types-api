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

func (h *PopulationTypes) Get(w http.ResponseWriter, req *http.Request) {
	body := struct {
		Items []interface{} `json:"items"`
	}{
		Items: make([]interface{}, 0),
	}
	responder := h.responder
	ctx := req.Context()

	responder.JSON(ctx, w, 200, body)
}
