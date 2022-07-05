package handler

import (
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	dperrors "github.com/ONSdigital/dp-api-clients-go/v2/errors"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/contract"
	"github.com/go-chi/chi/v5"

	dprequest "github.com/ONSdigital/dp-net/v2/request"
	"github.com/pkg/errors"
)

// Areas handles requests to /area-types
type Areas struct {
	cfg      *config.Config
	datasets datasetAPIClient
	respond  responder
	ctblr    cantabularClient
}

// NewAreas returns a new Areas handler
func NewAreas(cfg *config.Config, d datasetAPIClient, r responder, c cantabularClient) *Areas {
	return &Areas{
		cfg:      cfg,
		datasets: d,
		respond:  r,
		ctblr:    c,
	}
}

func (a *Areas) authenticate(r *http.Request) bool {
	var authorised bool

	var hasCallerIdentity, hasUserIdentity bool
	callerIdentity := dprequest.Caller(r.Context())

	if callerIdentity != "" {
		hasCallerIdentity = true
	}

	userIdentity := dprequest.User(r.Context())
	if userIdentity != "" {
		hasUserIdentity = true
	}

	if hasCallerIdentity || hasUserIdentity {
		authorised = true
	}

	return authorised
}

// Get is the handler for GET /population-types/{population-type}/area-types/{area-type}/areas
func (h *Areas) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	authToken := "" //will be passing empty auth token to GetDatasets for public call
	if h.cfg.EnablePrivateEndpoints {
		authToken = h.cfg.ServiceAuthToken
	}

	datasetName := chi.URLParam(r, "population-type")
	areaType := chi.URLParam(r, "area-type")
	category := r.URL.Query().Get("q")

	areaTypeReq := cantabular.GetAreasRequest{
		Dataset:  datasetName,
		Variable: areaType,
		Category: category,
	}

	datasets, err := h.datasets.GetDatasets(
		ctx,
		"",
		authToken,
		"",
		&dataset.QueryParams{
			IsBasedOn: datasetName,
			Limit:     1000,
		},
	)
	if err != nil {
		h.respond.Error(
			ctx,
			w,
			dperrors.StatusCode(err),
			errors.New("failed to get area types: failed to get areas"),
		)
		return

	}

	if datasets.Count == 0 {
		h.respond.Error(
			ctx,
			w,
			http.StatusNotFound,
			errors.New("areas not found"),
		)
		return
	}

	//if dataset is not empty, for private endpoint do extra checks
	if h.cfg.EnablePrivateEndpoints {
		var isPublished bool
		for _, d := range datasets.Items {
			if d.Current != nil {
				isPublished = true
				break
			}
		}

		if !isPublished {
			h.respond.Error(
				ctx,
				w,
				http.StatusNotFound,
				errors.New("areas not found for private call"),
			)
			return
		}
	}

	areas, err := h.ctblr.GetAreas(ctx, areaTypeReq)
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

	h.respond.JSON(ctx, w, http.StatusOK, toAreasResponse(areas))
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
