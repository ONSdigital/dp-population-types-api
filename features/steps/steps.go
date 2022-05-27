package steps

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/cucumber/godog"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular/gql"
)

func (c *PopulationTypesComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`^a list of named cantabular population types is returned$`, c.aListOfNamedCantabularPopulationTypesIsReturned)
	ctx.Step(`^cantabular is unresponsive$`, c.cantabularIsUnresponsive)
	ctx.Step(`^I access the root population types endpoint$`, c.iAccessTheRootPopulationTypesEndpoint)
	ctx.Step(`^I have some population types in cantabular$`, c.iHaveSomePopulationTypesInCantabular)
	ctx.Step(`^the service responds with "([^"]*)" http code and an internal server error saying "([^"]*)"$`, c.theServiceRespondsWithAnInternalServerErrorSaying)
	ctx.Step(`^a geography query response is available from Cantabular api extension$`, c.theFollowingCantabularResponseIsAvailable)
	ctx.Step(`^an error json response is returned from Cantabular api extension$`, c.anErrorJsonResponseIsReturnedFromCantabularApiExtension)
	ctx.Step(`^a list of area types is returned$`, c.aListOfAreaTypesIsReturned)
}

func (c *PopulationTypesComponent) aListOfNamedCantabularPopulationTypesIsReturned() error {
	return c.apiFeature.IShouldReceiveTheFollowingJSONResponse(&godog.DocString{
		Content: `{"items":[{"name": "dataset 1"}, {"name": "dataset 2"}, {"name": "dataset 3"}]}`,
	})
}

func (c *PopulationTypesComponent) iAccessTheRootPopulationTypesEndpoint() error {
	return c.apiFeature.IGet("/population-types")
}

func (c *PopulationTypesComponent) iHaveSomePopulationTypesInCantabular() error {
	c.fakeCantabularDatasets = []string{"dataset 1", "dataset 2", "dataset 3"}
	return nil
}

func (c *PopulationTypesComponent) cantabularIsUnresponsive() error {
	c.fakeCantabularIsUnresponsive = true
	return nil
}

func (c *PopulationTypesComponent) theServiceRespondsWithAnInternalServerErrorSaying(expectedHttpCode int, expected string) error {
	resp := c.apiFeature.HttpResponse
	if resp.StatusCode != expectedHttpCode {
		return fmt.Errorf("expected: %d. actual: %d", http.StatusInternalServerError, resp.StatusCode)
	}
	body, err := ioutil.ReadAll(c.apiFeature.HttpResponse.Body)
	if err != nil {
		return err
	}
	if !strings.Contains(string(body), expected) {
		return fmt.Errorf("expected to contain: %s. actual: %s", fakeCantabularFailedToRespondErrorMessage, string(body))
	}
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

func (c *PopulationTypesComponent) aListOfAreaTypesIsReturned() error {
	return c.apiFeature.IShouldReceiveTheFollowingJSONResponse(&godog.DocString{
		Content: `
			  {
				"area-types":[
				  {
					"id":"country",
					"label":"Country",
					"total_count": 2
				  },
				  {
					"id":"city",
					"label":"City",
					"total_count": 3
				  }
				]
			  }
			`,
	})
}
