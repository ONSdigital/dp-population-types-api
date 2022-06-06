Feature: Area Types

  Background:
    Given private endpoints are not enabled

  Scenario: Getting published area types
    Given the following response is returned from dp-dataset-api:
    """
    {
    "total_count": 1
    }
    """

    And a geography query response is available from Cantabular api extension

    When I GET "/population-types/Example/area-types"

    Then the HTTP status code should be "200"
    And I should receive the following JSON response:
    """
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
    """

  Scenario: Getting unpublished area types
    Given the following response is returned from dp-dataset-api:
    """
    {
    "total_count": 0
    }
    """

    And a geography query response is available from Cantabular api extension

    When I GET "/population-types/Example/area-types"
    Then the HTTP status code should be "404"

  Scenario: Getting area-types not found
    When an error json response is returned from Cantabular api extension

    And I GET "/population-types/Inexistent/area-types"
    Then an internal server error saying "failed to get area-types: error(s) returned by graphQL query" is returned
    And the HTTP status code should be "404"
