package steps

import (
	"github.com/cucumber/godog"
)

func (c *PopulationTypesComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	c.apiFeature.RegisterSteps(ctx)

	ctx.Step(`^a list of named cantabular population types is returned$`, c.aListOfNamedCantabularPopulationTypesIsReturned)
	ctx.Step(`^cantabular is unresponsive$`, c.cantabularIsUnresponsive)
	ctx.Step(`^I access the root population types endpoint$`, c.iAccessTheRootPopulationTypesEndpoint)
	ctx.Step(`^I have some population types in cantabular$`, c.iHaveSomePopulationTypesInCantabular)
	ctx.Step(`^the service responds with an internal server error saying "([^"]*)"$`, c.theServiceRespondsWithAnInternalServerErrorSaying)
}

func (c *PopulationTypesComponent) aListOfNamedCantabularPopulationTypesIsReturned() error {
	return c.apiFeature.IShouldReceiveTheFollowingJSONResponseWithStatus(
		"200",
		&godog.DocString{Content: `{ 
			"items": [
				{ "name": "dataset 1" },
				{ "name": "dataset 2" },
				{ "name": "dataset 3" }
			]
		}`},
	)
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

func (c *PopulationTypesComponent) theServiceRespondsWithAnInternalServerErrorSaying(arg1 string) error {
	return c.apiFeature.IShouldReceiveTheFollowingResponse(
		&godog.DocString{MediaType: "text/plain", Content: "failed to fetch population types"})
}
