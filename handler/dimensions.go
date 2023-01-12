package handler

import (
	"context"
	"net/http"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/contract"
	"github.com/ONSdigital/log.go/v2/log"
	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"
)

type Dimensions struct {
	cfg        *config.Config
	respond    responder
	cantabular cantabularClient
	datasets   datasetAPIClient
}

// NewDimensions returns a new AreaTypes handler set
func NewDimensions(cfg *config.Config, r responder, c cantabularClient, d datasetAPIClient) *Dimensions {
	return &Dimensions{
		cfg:        cfg,
		respond:    r,
		cantabular: c,
		datasets:   d,
	}
}

// GetAll is the handler for GET /population-types/{population-type}/dimensions
func (h *Dimensions) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := contract.GetDimensionsRequest{
		PopulationType: chi.URLParam(r, "population-type"),
	}
	if err := parseRequest(r, &req); err != nil {
		h.respond.Error(ctx, w, http.StatusBadRequest, &Error{
			err: errors.Wrap(err, "invalid request"),
		})
		return
	}

	logData := log.Data{
		"request": req,
	}

	// only return results for published population-types on web
	if !h.cfg.EnablePrivateEndpoints {
		if err := h.published(ctx, req.PopulationType); err != nil {
			h.respond.Error(ctx, w, http.StatusNotFound, &Error{
				err:     errors.Wrap(err, "failed to check published state"),
				message: "population type not found",
				logData: logData,
			})
			return
		}
	}

	cReq := cantabular.GetDimensionsRequest{
		Dataset: req.PopulationType,
		Text:    req.SearchText,
		PaginationParams: cantabular.PaginationParams{
			Limit:  req.QueryParams.Limit,
			Offset: req.QueryParams.Offset,
		},
	}

	res, err := h.cantabular.GetDimensions(ctx, cReq)
	if err != nil {
		h.respond.Error(ctx, w, h.cantabular.StatusCode(err), &Error{
			err:     errors.Wrap(err, "failed to get dimensions"),
			message: "failed to get dimensions",
			logData: logData,
		})
		return
	}

	resp := contract.GetDimensionsResponse{
		PaginationResponse: contract.PaginationResponse{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}

	if res != nil {
		resp.Count = res.Count
		resp.TotalCount = res.TotalCount
		for _, edge := range res.Dataset.Variables.Search.Edges {
			resp.Dimensions = append(resp.Dimensions, contract.Dimension{
				ID:          edge.Node.Name,
				Label:       edge.Node.Label,
				Description: edge.Node.Description,
				TotalCount:  edge.Node.Categories.TotalCount,
			})
		}
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}

// GetDescription is the handler for GET /population-types/{population-type}/dimensions/description
func (h *Dimensions) GetDescription(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := contract.GetDimensionsDescriptionRequest{
		PopulationType: chi.URLParam(r, "population-type"),
	}
	if err := parseRequest(r, &req); err != nil {
		h.respond.Error(ctx, w, http.StatusBadRequest, &Error{
			err: errors.Wrap(err, "invalid request"),
		})
		return
	}

	logData := log.Data{
		"request": req,
	}

	// only return results for published population-types on web
	if !h.cfg.EnablePrivateEndpoints {
		if err := h.published(ctx, req.PopulationType); err != nil {
			h.respond.Error(ctx, w, http.StatusNotFound, &Error{
				err:     errors.Wrap(err, "failed to check published state"),
				message: "population type not found",
				logData: logData,
			})
			return
		}
	}

	cReq := cantabular.GetDimensionsDescriptionRequest{
		Dataset:        req.PopulationType,
		DimensionNames: req.DimensionNames,
	}

	res, err := h.cantabular.GetDimensionsDescription(ctx, cReq)
	if err != nil {
		h.respond.Error(ctx, w, h.cantabular.StatusCode(err), &Error{
			err:     errors.Wrap(err, "failed to get dimensions"),
			message: "failed to get dimensions",
			logData: logData,
		})
		return
	}

	resp := contract.GetDimensionsResponse{
		PaginationResponse: contract.PaginationResponse{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}

	if res != nil {
		resp.Count = res.Count
		resp.TotalCount = res.TotalCount
		for _, edge := range res.Dataset.Variables.Edges {
			resp.Dimensions = append(resp.Dimensions, contract.Dimension{
				ID:          edge.Node.Name,
				Label:       edge.Node.Label,
				Description: edge.Node.Description,
				TotalCount:  edge.Node.Categories.TotalCount,
			})
		}
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}

// GetCategorisations is the handler for GET /population-types/{population-type}/dimensions/{dimension}/categoristations
func (h *Dimensions) GetCategorisations(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := contract.GetCategorisationsRequest{
		PopulationType: chi.URLParam(r, "population-type"),
		Variable:       chi.URLParam(r, "dimension"),
	}

	logData := log.Data{
		"population_type": req.PopulationType,
		"dimension":       req.Variable,
	}

	if err := parseRequest(r, &req); err != nil {
		h.respond.Error(ctx, w, http.StatusBadRequest, &Error{
			err: errors.Wrap(err, "invalid request"),
		})
		return
	}

	// only return results for published population-types on web
	if !h.cfg.EnablePrivateEndpoints {
		if err := h.published(ctx, req.PopulationType); err != nil {
			h.respond.Error(ctx, w, http.StatusNotFound, &Error{
				err:     errors.Wrap(err, "failed to check published state"),
				message: "population type not found",
				logData: logData,
			})
			return
		}
	}
	cReq := cantabular.GetCategorisationsRequest{
		Dataset:  req.PopulationType,
		Variable: req.Variable,
		PaginationParams: cantabular.PaginationParams{
			Limit:  req.QueryParams.Limit,
			Offset: req.QueryParams.Offset,
		},
	}

	res, err := h.cantabular.GetCategorisations(ctx, cReq)
	if err != nil {
		h.respond.Error(ctx, w, h.cantabular.StatusCode(err), &Error{
			err:     errors.Wrap(err, "failed to get categorisations"),
			message: "failed to get categorisations",
			logData: logData,
		})
		return
	}

	resp := contract.GetCategorisationsResponse{
		PaginationResponse: contract.PaginationResponse{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}

	if res != nil {
		resp.Count = res.Count
		resp.TotalCount = res.TotalCount
		for _, edge := range res.Dataset.Variables.Edges {
			// Check if the mapFrom array is populated, if so it is not the base variable and the results will be here
			if len(edge.Node.MapFrom) > 0 {
				for _, mapFrom := range edge.Node.MapFrom {

					for _, mapEdge := range mapFrom.Edges {
						resp.TotalCount = mapEdge.Node.IsSourceOf.TotalCount
						for _, isSourceOf := range mapEdge.Node.IsSourceOf.Edges {

							cats := []contract.DimensionCategory{}
							for _, categories := range isSourceOf.Node.Categories.Edges {
								cats = append(cats, contract.DimensionCategory{
									ID:    categories.Node.Code,
									Label: categories.Node.Label,
								})
							}

							resp.Items = append(resp.Items, contract.Category{
								Id:         isSourceOf.Node.Name,
								Label:      isSourceOf.Node.Label,
								Categories: cats,
							})
						}
					}
				}
			} else {
				// This is the base variable that is queried so the categorisations will be in the IsSourceOfArray
				resp.TotalCount = edge.Node.IsSourceOf.TotalCount
				for _, sourceOf := range edge.Node.IsSourceOf.Edges {

					cats := []contract.DimensionCategory{}
					for _, categories := range sourceOf.Node.Categories.Edges {
						cats = append(cats, contract.DimensionCategory{
							ID:    categories.Node.Code,
							Label: categories.Node.Label,
						})
					}

					resp.Items = append(resp.Items, contract.Category{
						Id:         sourceOf.Node.Name,
						Label:      sourceOf.Node.Label,
						Categories: cats,
					})
				}
			}

		}
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}

// GetBase returns the base variables for a given categorisation
func (h *Dimensions) GetBaseVariable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := contract.GetBaseVariableRequest{
		PopulationType: chi.URLParam(r, "population-type"),
		Variable:       chi.URLParam(r, "dimension"),
	}

	logData := log.Data{
		"population_type": req.PopulationType,
		"dimension":       req.Variable,
	}

	if err := parseRequest(r, &req); err != nil {
		h.respond.Error(ctx, w, http.StatusBadRequest, &Error{
			err: errors.Wrap(err, "invalid request"),
		})
		return
	}

	// only return results for published population-types on web
	if !h.cfg.EnablePrivateEndpoints {
		if err := h.published(ctx, req.PopulationType); err != nil {
			h.respond.Error(ctx, w, http.StatusNotFound, &Error{
				err:     errors.Wrap(err, "failed to check published state"),
				message: "base variable not found",
				logData: logData,
			})
			return
		}
	}
	cReq := cantabular.GetBaseVariableRequest{
		Dataset:  req.PopulationType,
		Variable: req.Variable,
	}

	res, err := h.cantabular.GetBaseVariable(ctx, cReq)
	if err != nil {
		h.respond.Error(ctx, w, h.cantabular.StatusCode(err), &Error{
			err:     errors.Wrap(err, "failed to get base variable"),
			message: "failed to get base variable",
			logData: logData,
		})
		return
	}

	resp := contract.GetBaseVariableResponse{}

	if res != nil {
		if len(res.Dataset.Variables.Edges) == 0 ||
			len(res.Dataset.Variables.Edges[0].Node.MapFrom) == 0 ||
			len(res.Dataset.Variables.Edges[0].Node.MapFrom[0].Edges) == 0 {
			h.respond.Error(ctx, w, http.StatusInternalServerError, &Error{
				err:     errors.Wrap(err, "cantabular returned unexpected empty list"),
				message: "failed to get base variable",
				logData: logData,
			})
			return
		}
		resp.Name = res.Dataset.Variables.Edges[0].Node.MapFrom[0].Edges[0].Node.Name
		resp.Label = res.Dataset.Variables.Edges[0].Node.MapFrom[0].Edges[0].Node.Label
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)

}

func (h *Dimensions) published(ctx context.Context, populationType string) error {
	datasets, err := h.datasets.GetDatasets(
		ctx,
		"",
		"",
		"",
		&dataset.QueryParams{
			IsBasedOn: populationType,
			Limit:     1000,
		},
	)
	if err != nil {
		return errors.Wrap(err, "failed to get datasets for population type")
	}

	if datasets.TotalCount == 0 {
		return errors.New("no published datasets found for population type")
	}

	return nil
}
