package handler

import (
	"context"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/pkg/errors"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/dataset"
	"github.com/ONSdigital/dp-population-types-api/config"
	"github.com/ONSdigital/dp-population-types-api/contract"
	"github.com/ONSdigital/log.go/v2/log"
)

// AreaTypes handles requests sent to the /population-types/{population-type}/area-types
// set of routes
type AreaTypes struct {
	cfg        *config.Config
	respond    responder
	cantabular cantabularClient
	datasets   datasetAPIClient
}

// NewAreaTypes returns a new AreaTypes handler set
func NewAreaTypes(cfg *config.Config, r responder, c cantabularClient, d datasetAPIClient) *AreaTypes {
	return &AreaTypes{
		cfg:        cfg,
		respond:    r,
		cantabular: c,
		datasets:   d,
	}
}

// Get is the handler for GET /population-types/{population-type}/area-types
func (h *AreaTypes) Get(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := contract.GetAreaTypesRequest{
		PopulationType: chi.URLParam(r, "population-type"),
	}

	if err := parseRequest(r, &req); err != nil {
		h.respond.Error(ctx, w, http.StatusBadRequest, &Error{
			err: errors.Wrap(err, "invalid request"),
		})
		return
	}

	cReq := cantabular.GetGeographyDimensionsRequest{
		Dataset: req.PopulationType,
		PaginationParams: cantabular.PaginationParams{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}

	logData := log.Data{
		"population_type": cReq.Dataset,
	}

	// only return results for published population-types on web
	if !h.cfg.EnablePrivateEndpoints {
		if err := h.published(ctx, cReq.Dataset); err != nil {
			h.respond.Error(ctx, w, http.StatusNotFound, &Error{
				err:     errors.Wrap(err, "failed to check published state"),
				message: "population type not found",
				logData: logData,
			})
			return
		}
	}

	res, err := h.cantabular.GetGeographyDimensions(ctx, cReq)
	if err != nil {
		h.respond.Error(ctx, w, h.cantabular.StatusCode(err), &Error{
			err:     errors.Wrap(err, "failed to get geography dimensions"),
			message: "failed to get area-types",
			logData: logData,
		})
		return
	}

	resp := contract.GetAreaTypesResponse{
		PaginationResponse: contract.PaginationResponse{
			Limit:  cReq.Limit,
			Offset: cReq.Offset,
		},
	}

	if res != nil {
		resp.Count = res.Count
		resp.TotalCount = res.TotalCount

		for _, edge := range res.Dataset.Variables.Edges {

			i := 0
			if len(edge.Node.Meta.ONSVariable.GeographyHierarchyOrder) > 0 {
				i, err = strconv.Atoi(edge.Node.Meta.ONSVariable.GeographyHierarchyOrder)
				if err != nil {
					h.respond.Error(ctx, w, http.StatusBadRequest, &Error{
						err: errors.Wrap(err, "unable to cast geography order field to int"),
					})
					return
				}
			}

			resp.AreaTypes = append(resp.AreaTypes, contract.AreaType{
				ID:                      edge.Node.Name,
				Label:                   edge.Node.Label,
				Description:             edge.Node.Description,
				TotalCount:              edge.Node.Categories.TotalCount,
				GeographyHierarchyOrder: i,
			})
		}
		sort.Slice(resp.AreaTypes, func(i, j int) bool {
			return resp.AreaTypes[i].GeographyHierarchyOrder > resp.AreaTypes[j].GeographyHierarchyOrder
		})
	}

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}

// GetParents is the handler for /population-types/{population-type}/area-types/{area-type}/parents
func (h *AreaTypes) GetParents(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req := contract.GetAreaTypeParentsRequest{
		PopulationType: chi.URLParam(r, "population-type"),
		AreaType:       chi.URLParam(r, "area-type"),
	}

	if err := parseRequest(r, &req); err != nil {
		h.respond.Error(ctx, w, http.StatusBadRequest, &Error{
			err: errors.Wrap(err, "invalid request"),
		})
		return
	}

	cReq := cantabular.GetParentsRequest{
		PaginationParams: cantabular.PaginationParams{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
		Dataset:  req.PopulationType,
		Variable: req.AreaType,
	}

	logData := log.Data{
		"population_type":   cReq.Dataset,
		"area_type":         cReq.Variable,
		"pagination_params": cReq.PaginationParams,
	}
	// only return results for published population-types on web
	if !h.cfg.EnablePrivateEndpoints {
		if err := h.published(ctx, cReq.Dataset); err != nil {
			h.respond.Error(ctx, w, http.StatusNotFound, &Error{
				err:     errors.Wrap(err, "failed to check published state"),
				message: "population type not found",
				logData: logData,
			})
			return
		}
	}

	res, err := h.cantabular.GetParents(ctx, cReq)
	if err != nil {
		h.respond.Error(ctx, w, h.cantabular.StatusCode(err), &Error{
			err:     errors.Wrap(err, "failed to get parents"),
			message: "failed to get parents",
			logData: logData,
		})
		return
	}

	resp := contract.GetAreaTypeParentsResponse{
		PaginationResponse: contract.PaginationResponse{
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	}

	if l := len(res.Dataset.Variables.Edges); l != 1 {
		logData["edges_expected_length"] = 1
		logData["edges_length"] = l
		h.respond.Error(ctx, w, http.StatusInternalServerError, &Error{
			err:     errors.New("invalid response from Cantabular"),
			message: "failed to get parents",
			logData: logData,
		})
		return
	}

	resp.Count = res.Count
	resp.TotalCount = res.TotalCount
	for _, e := range res.Dataset.Variables.Edges[0].Node.IsSourceOf.Edges {

		i := 0
		if len(e.Node.Meta.ONSVariable.GeographyHierarchyOrder) > 0 {
			i, err = strconv.Atoi(e.Node.Meta.ONSVariable.GeographyHierarchyOrder)
			if err != nil {
				h.respond.Error(ctx, w, http.StatusBadRequest, &Error{
					err: errors.Wrap(err, "unable to cast geography order field to int"),
				})
				return
			}
		}

		resp.AreaTypes = append(resp.AreaTypes, contract.AreaType{
			ID:                      e.Node.Name,
			Label:                   e.Node.Label,
			TotalCount:              e.Node.Categories.TotalCount,
			GeographyHierarchyOrder: i,
		})
	}

	sort.Slice(resp.AreaTypes, func(i, j int) bool {
		return resp.AreaTypes[i].GeographyHierarchyOrder > resp.AreaTypes[j].GeographyHierarchyOrder
	})

	h.respond.JSON(ctx, w, http.StatusOK, resp)
}

// GetParentAreaCount is the handler for /population-types/{population-type}/area-types/{area-type}/parents
func (h *AreaTypes) GetParentAreaCount(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	cReq := cantabular.GetParentAreaCountRequest{
		Dataset:   chi.URLParam(r, "population-type"),
		Variable:  chi.URLParam(r, "area-type"),
		Parent:    chi.URLParam(r, "parent-area-type"),
		SVariable: r.URL.Query().Get("svar"),
		Codes:     strings.Split(r.URL.Query().Get("areas"), ","),
	}

	logData := log.Data{
		"population_type":  cReq.Dataset,
		"area_type":        cReq.Variable,
		"parent_area_type": cReq.Parent,
		"codes":            cReq.Codes,
	}

	// only return results for published population-types on web
	if !h.cfg.EnablePrivateEndpoints {
		if err := h.published(ctx, cReq.Dataset); err != nil {
			h.respond.Error(ctx, w, http.StatusNotFound, &Error{
				err:     errors.Wrap(err, "failed to check published state"),
				message: "population type not found",
				logData: logData,
			})
			return
		}
	}

	res, err := h.cantabular.GetParentAreaCount(ctx, cReq)
	if err != nil {
		h.respond.Error(ctx, w, h.cantabular.StatusCode(err), &Error{
			err:     errors.Wrap(err, "failed to get parent areas count"),
			message: "failed to get parent areas count",
			logData: logData,
		})
		return
	}

	h.respond.JSON(ctx, w, http.StatusOK, res.Dimension.Count)
}

func (h *AreaTypes) published(ctx context.Context, populationType string) error {
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
