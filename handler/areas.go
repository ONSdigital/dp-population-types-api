package handler

import (
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-population-types-api/contract"

	"github.com/pkg/errors"
)

// Areas handles requests to /area-types
type Areas struct {
	respond responder
	ctblr   cantabularClient
}

// NewCoverage returns a new Areas handler
func NewCoverage(r responder, c cantabularClient) *Areas {
	return &Areas{
		respond: r,
		ctblr:   c,
	}
}

// Get is the handler for GET /coverage
func (h *Areas) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var req contract.GetAreasRequest
	if err := parseRequest(r, &req); err != nil {
		h.respond.Error(
			ctx,
			w,
			http.StatusBadRequest,
			fmt.Errorf("failed to parse request: %w", err),
		)
		return
	}

	areaTypeReq := cantabular.GetAreasRequest{
		Dataset:  req.Dataset,
		Variable: req.AreaType,
		Category: req.Text,
	}

	coverage, err := h.ctblr.GetAreas(ctx, areaTypeReq)
	if err != nil {
		msg := "failed to get areas"
		h.respond.Error(
			ctx,
			w,
			h.ctblr.StatusCode(err),
			&Error{
				err:     errors.Wrap(err, msg),
				message: msg,
			},
		)
		return
	}

	h.respond.JSON(ctx, w, http.StatusOK, toAreasResponse(coverage))
}

// toAreasResponse converts a cantabular.GetAreasResponse to a flattened contract.GetAreasResponse.
func toAreasResponse(res *cantabular.GetAreasResponse) contract.GetAreasResponse {
	var resp contract.GetAreasResponse

	for _, variable := range res.Dataset.RuleBase.IsSourceOf.Search.Edges {
		for _, category := range variable.Node.Categories.Search.Edges {
			resp.Areas = append(resp.Areas, contract.Areas{
				ID:       category.Node.Code,
				Label:    category.Node.Label,
				AreaType: variable.Node.Name,
			})
		}
	}

	return resp
}
