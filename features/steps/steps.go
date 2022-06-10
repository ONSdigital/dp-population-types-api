package steps

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cucumber/godog"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
)

func (c *PopulationTypesComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`^cantabular is unresponsive$`, c.cantabularIsUnresponsive)
	ctx.Step(`^a geography query response is available from Cantabular api extension$`, c.theFollowingCantabularResponseIsAvailable)
	ctx.Step(`^an error json response is returned from Cantabular api extension$`, c.anErrorJsonResponseIsReturnedFromCantabularApiExtension)
	ctx.Step(`^I have the following population types in cantabular$`, c.iHaveTheFollowingPopulationTypesInCantabular)
	ctx.Step(`^the following datasets based on "([^"]*)" are available$`, c.theFollowingDatasetsBasedOnAreAvailable)
}

func (c *PopulationTypesComponent) iHaveTheFollowingPopulationTypesInCantabular(body *godog.DocString) error {
	var populationTypes []string
	if err := json.Unmarshal([]byte(body.Content), &populationTypes); err != nil {
		return fmt.Errorf("failed to unmarshal population types: %w", err)
	}
	c.fakeCantabularDatasets = populationTypes
	return nil
}

func (c *PopulationTypesComponent) iHaveTheFollowingDatasetsFromDatasetAPI(body *godog.DocString) error {
	var populationTypes []string
	if err := json.Unmarshal([]byte(body.Content), &populationTypes); err != nil {
		return fmt.Errorf("failed to unmarshal population types: %w", err)
	}
	c.fakeCantabularDatasets = populationTypes
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

// theFollowingCantabularResponseIsAvailable generates a mocked response for Cantabular Server POST /graphql
func (c *PopulationTypesComponent) theFollowingCantabularResponseIsAvailable() error {
	gd := &cantabular.GetGeographyDimensionsResponse{Dataset: gql.DatasetRuleBase{RuleBase: gql.RuleBase{
		IsSourceOf: gql.Variables{Edges: []gql.Edge{
			{
				Node: gql.Node{
					Name:       "country",
					Label:      "Country",
					Categories: gql.Categories{TotalCount: 2},
					MapFrom: []gql.Variables{{Edges: []gql.Edge{{
						Node: gql.Node{
							Name:       "city",
							Label:      "City",
							FilterOnly: "false",
						},
					}}}},
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
		}},
	}}}
	c.fakeCantabularGeoDimensions = gd
	return nil
}

func (c *PopulationTypesComponent) anErrorJsonResponseIsReturnedFromCantabularApiExtension() error {
	c.fakeCantabularGeoDimensions = nil
	c.fakeCantabularIsUnresponsive = true
	return nil
}
