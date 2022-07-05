package steps

import (
	"context"
	"fmt"
	"net/http"

	"github.com/cucumber/godog"
	"github.com/pkg/errors"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"

	"github.com/ONSdigital/log.go/v2/log"
)

func (c *PopulationTypesComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^private endpoints are enabled`, c.privateEndpointsAreEnabled)
	ctx.Step(`^private endpoints are not enabled`, c.privateEndpointsAreNotEnabled)
	ctx.Step(`^cantabular is unresponsive$`, c.cantabularIsUnresponsive)
	ctx.Step(`^a geography query response is available from Cantabular api extension$`, c.theFollowingCantabularResponseIsAvailable)
	ctx.Step(`^an error json response is returned from Cantabular api extension$`, c.anErrorJsonResponseIsReturnedFromCantabularApiExtension)
	ctx.Step(`^I have the following list datasets response available in cantabular$`, c.iHaveTheFollowingPopulationTypesInCantabular)
	ctx.Step(`^the following datasets based on "([^"]*)" are available$`, c.theFollowingDatasetsBasedOnAreAvailable)
	ctx.Step("^the dp-dataset-api is returning errors", c.datasetClientReturnsErrors)
	ctx.Step(`^cantabular api extension is healthy`, c.cantabularAPIExtIsHealthy)
	ctx.Step(`^cantabular server is healthy`, c.cantabularServerIsHealthy)
	ctx.Step(`^the following area query response is available from Cantabular api extension for the dataset "([^"]*)":$`, c.theFollowingCantabularAreaResponseIsAvailable)
	ctx.Step(`^the following area query response is available from Cantabular api extension for the dataset "([^"]*)", area type "([^"]*)" and text "([^"]*)":$`, c.theFollowingCantabularFilteredAreaResponseIsAvailable)
	ctx.Step(`^the following geography response is available from Cantabular for the dataset "([^"]*)":$`, c.theFollowingCantabularGeographyResponseIsAvailable)
}

func (c *PopulationTypesComponent) iHaveTheFollowingPopulationTypesInCantabular(body *godog.DocString) error {
	data := cantabular.QueryData{}

	b, err := data.Encode(cantabular.QueryListDatasets)
	if err != nil {
		return errors.Wrap(err, "failed to encode GraphQL query")
	}

	c.CantabularApiExt.NewHandler().
		Post("/graphql").
		AssertBody(b.Bytes()).
		Reply(http.StatusOK).
		BodyString(body.Content)
	return nil
}

func (c *PopulationTypesComponent) datasetClientReturnsErrors() {
	url := fmt.Sprintf(
		`/datasets?offset=0&limit=100&is_based_on=%s`,
		"Example",
	)

	c.datasetAPI.NewHandler().
		Get(url).
		Reply(http.StatusInternalServerError).
		BodyString("some test error")
}

func (c *PopulationTypesComponent) privateEndpointsAreNotEnabled() error {
	c.Config.EnablePrivateEndpoints = false
	return nil
}

func (c *PopulationTypesComponent) privateEndpointsAreEnabled() error {
	c.Config.EnablePrivateEndpoints = true
	return nil
}

func (c *PopulationTypesComponent) cantabularIsUnresponsive() error {
	c.fakeCantabularIsUnresponsive = true
	return nil
}

func (c *PopulationTypesComponent) theFollowingDatasetsBasedOnAreAvailable(populationType string, body *godog.DocString) error {
	url := "/datasets?offset=0&limit=1000&is_based_on=" + populationType
	c.datasetAPI.NewHandler().
		Get(url).
		Reply(http.StatusOK).
		BodyString(body.Content)

	return nil
}

func (c *PopulationTypesComponent) theFollowingCantabularGeographyResponseIsAvailable(dataset string, body *godog.DocString) error {
	data := cantabular.QueryData{
		Dataset: dataset,
	}

	b, err := data.Encode(cantabular.QueryGeographyDimensions)
	if err != nil {
		return errors.Wrap(err, "failed to encode GraphQL query")
	}
	log.Info(context.Background(), "DEBUG", log.Data{"DEBUG": string(b.Bytes())})
	// create graphql handler with expected query body
	c.CantabularApiExt.NewHandler().
		Post("/graphql").
		AssertBody(b.Bytes()).
		Reply(http.StatusOK).
		BodyString(body.Content)

	return nil
}

// theFollowingCantabularResponseIsAvailable generates a mocked response for Cantabular Server POST /graphql
func (c *PopulationTypesComponent) theFollowingCantabularResponseIsAvailable() error {
	/*gd := &cantabular.GetGeographyDimensionsResponse{
		Dataset: gql.DatasetRuleBase{
			RuleBase: gql.RuleBase{
				IsSourceOf: gql.Variables{
					Edges: []gql.Edge{
						{
							Node: gql.Node{
								Name:       "country",
								Label:      "Country",
								Categories: gql.Categories{TotalCount: 2},
								MapFrom: []gql.Variables{
									{
										Edges: []gql.Edge{
											{
												Node: gql.Node{
													Name:       "city",
													Label:      "City",
													FilterOnly: "false",
												},
											},
										},
									},
								},
							},
						},
						{
							Node: gql.Node{
								Name:       "city",
								Label:      "City",
								Categories: gql.Categories{TotalCount: 3},
								MapFrom:    []gql.Variables{},
							},
						},
					},
				},
			},
		},
	}*/
	//c.fakeCantabularGeoDimensions = gd
	return nil
}

func (c *PopulationTypesComponent) anErrorJsonResponseIsReturnedFromCantabularApiExtension() error {
	//c.fakeCantabularGeoDimensions = nil
	c.fakeCantabularIsUnresponsive = true
	return nil
}

// cantabularAPIExtIsHealthy generates a mocked healthy response for cantabular server
func (c *PopulationTypesComponent) cantabularAPIExtIsHealthy() error {
	const res = `{"status": "OK"}`
	c.CantabularApiExt.NewHandler().
		Get("/graphql?query={}").
		Reply(http.StatusOK).
		BodyString(res)
	return nil
}

// cantabularServerIsHealthy generates a mocked healthy response for cantabular server
func (c *PopulationTypesComponent) cantabularServerIsHealthy() error {
	const res = `{"status": "OK"}`
	c.CantabularSrv.NewHandler().
		Get("/v9/datasets").
		Reply(http.StatusOK).
		BodyString(res)
	return nil
}

func (c *PopulationTypesComponent) theFollowingCantabularAreaResponseIsAvailable(dataset string, cb *godog.DocString) error {
	data := cantabular.QueryData{
		Dataset: dataset,
	}

	b, err := data.Encode(cantabular.QueryAreas)
	if err != nil {
		return errors.Wrap(err, "failed to encode GraphQL query")
	}

	// create graphql handler with expected query body
	c.CantabularApiExt.NewHandler().
		Post("/graphql").
		AssertBody(b.Bytes()).
		Reply(http.StatusOK).
		BodyString(cb.Content)
	/*
		var resp = &struct {
			Data   cantabular.GetAreasResponse `json:"data"`
			Errors []gql.Error                 `json:"errors,omitempty"`
		}{}
		if err := json.Unmarshal([]byte(cb.Content), &resp); err != nil {
			return errors.Wrap(err, "failed to unmarshal GraphQL query")
		}
		c.fakeGetAreasResponse = &resp.Data*/
	return nil
}

func (c *PopulationTypesComponent) theFollowingCantabularFilteredAreaResponseIsAvailable(dataset, areaType, text string, cb *godog.DocString) error {
	data := cantabular.QueryData{
		Dataset:  dataset,
		Text:     areaType,
		Category: text,
	}

	b, err := data.Encode(cantabular.QueryAreas)
	if err != nil {
		return errors.Wrap(err, "failed to encode GraphQL query")
	}

	// create graphql handler with expected query body
	c.CantabularApiExt.NewHandler().
		Post("/graphql").
		AssertBody(b.Bytes()).
		Reply(http.StatusOK).
		BodyString(cb.Content)
	/*var resp = &struct {
		Data   cantabular.GetAreasResponse `json:"data"`
		Errors []gql.Error                 `json:"errors,omitempty"`
	}{}
	if err := json.Unmarshal([]byte(cb.Content), &resp); err != nil {
		return errors.Wrap(err, "failed to unmarshal GraphQL query")
	}
	c.fakeGetAreasResponse = &resp.Data*/
	return nil
}
