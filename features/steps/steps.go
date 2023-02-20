package steps

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cucumber/godog"

	"github.com/ONSdigital/dp-api-clients-go/v2/cantabular"
)

func (c *PopulationTypesComponent) RegisterSteps(ctx *godog.ScenarioContext) {
	ctx.Step(`^private endpoints are enabled`, c.privateEndpointsAreEnabled)
	ctx.Step(`^private endpoints are not enabled`, c.privateEndpointsAreNotEnabled)
	ctx.Step(`^cantabular is unresponsive$`, c.cantabularIsUnresponsive)
	ctx.Step(`^the following geography response is available from Cantabular:$`, c.theFollowingCantabularGeographyResponseIsAvailable)
	ctx.Step(`^the following dimensions response is available from Cantabular:$`, c.theFollowingCantabularDimensionsResponseIsAvailable)
	ctx.Step(`^the following dimensions description response is available from Cantabular:$`, c.theFollowingCantabularDimensionsDescriptionResponseIsAvailable)
	ctx.Step(`^I have the following population types in cantabular$`, c.iHaveTheFollowingPopulationTypesInCantabular)
	ctx.Step(`^the following datasets based on "([^"]*)" are available$`, c.theFollowingDatasetsBasedOnAreAvailable)
	ctx.Step(`^the dp-dataset-api is returning errors for datasets based on "([^"]*)"`, c.datasetClientReturnsErrors)
	ctx.Step(`^cantabular api extension is healthy`, c.cantabularAPIExtIsHealthy)
	ctx.Step(`^cantabular server is healthy`, c.cantabularServerIsHealthy)
	ctx.Step(`^the following GetArea query response is available from Cantabular:$`, c.theFollowingCantabularAreaResponseIsAvailable)
	ctx.Step(`^the following area query response is available from Cantabular:$`, c.theFollowingCantabularAreasResponseIsAvailable)
	ctx.Step(`^the following parents response is available from Cantabular:$`, c.theFollowingCantabularParentsResponseIsAvailable)
	ctx.Step(`^the following parents areas count response is available from Cantabular:$`, c.theFollowingCantabularParentAreaCountResponseIsAvailable)
	ctx.Step(`^the following blocked area response is available from Cantabular:$`, c.theFollowingBlockedAreaResponseIsAavailable)
	ctx.Step(`^the following categorisations response is available from Cantabular:$`, c.theFollowingCantabularCategorisationsResponseIsAvailable)
	ctx.Step(`^the cantabular area response is not found`, c.cantabularIsNotFound)
	ctx.Step(`^the cantabular area response is bad request`, c.cantabularIsBadRequest)
	ctx.Step(`^the cantabular response is bad gateway`, c.cantabularIsBadGateway)
	ctx.Step(`^the following base variable response is available from Cantabular:$`, c.theFollowingBaseVariableResponseIsAvailableFromCantabular)
	ctx.Step(`^the following dimension categories response is available from Cantabular:$`, c.theFollowingDimensionCategoryResponseIsAvailableFromCantabular)

}

func (c *PopulationTypesComponent) theFollowingDimensionCategoryResponseIsAvailableFromCantabular(body *godog.DocString) error {
	var response cantabular.GetDimensionCategoriesResponse
	if err := json.Unmarshal([]byte(body.Content), &response); err != nil {
		return fmt.Errorf("failed to unmarshal population types: %w", err)
	}
	c.fakeCantabular.GetDimensionCategoriesRespnse = &response
	return nil

}

func (c *PopulationTypesComponent) theFollowingBaseVariableResponseIsAvailableFromCantabular(body *godog.DocString) error {
	var response cantabular.GetBaseVariableResponse
	if err := json.Unmarshal([]byte(body.Content), &response); err != nil {
		return fmt.Errorf("failed to unmarshal population types: %w", err)
	}
	c.fakeCantabular.GetBaseVariableResponse = &response

	return nil

}

func (c *PopulationTypesComponent) iHaveTheFollowingPopulationTypesInCantabular(body *godog.DocString) error {
	var populationTypes cantabular.ListDatasetsResponse
	if err := json.Unmarshal([]byte(body.Content), &populationTypes); err != nil {
		return fmt.Errorf("failed to unmarshal population types: %w", err)
	}

	c.fakeCantabular.ListDatasetsResponse = &populationTypes
	return nil
}

func (c *PopulationTypesComponent) datasetClientReturnsErrors(populationType string) {
	url := fmt.Sprintf(
		`/datasets?offset=0&limit=1000&is_based_on=%s`,
		populationType,
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
	c.fakeCantabular.Healthy = false
	return nil
}

func (c *PopulationTypesComponent) cantabularIsNotFound() error {
	c.fakeCantabular.NotFound = true
	return nil
}

func (c *PopulationTypesComponent) cantabularIsBadRequest() error {
	c.fakeCantabular.BadRequest = true
	return nil
}

func (c *PopulationTypesComponent) cantabularIsBadGateway() error {
	c.fakeCantabular.BadGateway = true
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
func (c *PopulationTypesComponent) theFollowingCantabularGeographyResponseIsAvailable(body *godog.DocString) error {
	var resp cantabular.GetGeographyDimensionsResponse

	if err := json.Unmarshal([]byte(body.Content), &resp); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	c.fakeCantabular.GetGeographyDimensionsResponse = &resp
	return nil
}

// theFollowingCantabularResponseIsAvailable generates a mocked response for Cantabular Server POST /graphql
func (c *PopulationTypesComponent) theFollowingCantabularDimensionsResponseIsAvailable(body *godog.DocString) error {
	var resp cantabular.GetDimensionsResponse

	if err := json.Unmarshal([]byte(body.Content), &resp); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	c.fakeCantabular.GetDimensionsResponse = &resp
	return nil
}

func (c *PopulationTypesComponent) theFollowingCantabularDimensionsDescriptionResponseIsAvailable(body *godog.DocString) error {
	var resp cantabular.GetDimensionsResponse

	if err := json.Unmarshal([]byte(body.Content), &resp); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	c.fakeCantabular.GetDimensionsDescriptionResponse = &resp
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

func (c *PopulationTypesComponent) theFollowingCantabularAreasResponseIsAvailable(body *godog.DocString) error {
	var resp cantabular.GetAreasResponse

	if err := json.Unmarshal([]byte(body.Content), &resp); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	c.fakeCantabular.GetAreasResponse = &resp
	return nil
}

func (c *PopulationTypesComponent) theFollowingCantabularAreaResponseIsAvailable(body *godog.DocString) error {
	var resp cantabular.GetAreaResponse

	if err := json.Unmarshal([]byte(body.Content), &resp); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	c.fakeCantabular.GetAreaResponse = &resp
	return nil
}

func (c *PopulationTypesComponent) theFollowingCantabularParentsResponseIsAvailable(body *godog.DocString) error {
	var resp cantabular.GetParentsResponse

	if err := json.Unmarshal([]byte(body.Content), &resp); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	c.fakeCantabular.GetParentsResponse = &resp
	return nil
}

func (c *PopulationTypesComponent) theFollowingCantabularParentAreaCountResponseIsAvailable(body *godog.DocString) error {
	var resp cantabular.GetParentAreaCountResult
	resp.Dimension.Count = 1
	if err := json.Unmarshal([]byte(body.Content), &resp); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}
	c.fakeCantabular.GetParentAreaCountResult = &resp
	return nil
}

func (c *PopulationTypesComponent) theFollowingBlockedAreaResponseIsAavailable(body *godog.DocString) error {
	var resp cantabular.GetBlockedAreaCountResponse
	if err := json.Unmarshal([]byte(body.Content), &resp); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}
	c.fakeCantabular.GetBlockedAreaCountResult = &cantabular.GetBlockedAreaCountResult{
		Passed:     resp.Dataset.Table.Rules.Passed.Count,
		Total:      resp.Dataset.Table.Rules.Total.Count,
		Blocked:    resp.Dataset.Table.Rules.Blocked.Count,
		TableError: resp.Dataset.Table.Error,
	}
	return nil
}

func (c *PopulationTypesComponent) theFollowingCantabularCategorisationsResponseIsAvailable(body *godog.DocString) error {
	var resp cantabular.GetCategorisationsResponse

	if err := json.Unmarshal([]byte(body.Content), &resp); err != nil {
		return fmt.Errorf("failed to unmarshal body: %w", err)
	}

	c.fakeCantabular.GetCategorisationsResponse = &resp
	return nil
}
