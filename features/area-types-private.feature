Feature: Area Types
  Background:
    Given private endpoints are enabled
  Scenario: Getting Published area-types
    Given the following response is returned from dp-dataset-api:
    """
    {"total_count": 1}
    """
    And I have some population types in cantabular
    And a geography query response is available from Cantabular api extension
    And I am identified as "user@ons.gov.uk"
    And I am authorised

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

  Scenario: Getting Unpublished area-types


    Given I am identified as "user@ons.gov.uk"
    And I am authorised

    And I have some population types in cantabular
    And the following response is returned from dp-dataset-api:
    """
    {"total_count": 0}
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
