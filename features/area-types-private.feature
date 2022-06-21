Feature: Area Types
  Background:
    Given private endpoints are enabled

    And I am identified as "user@ons.gov.uk"

    And I am authorised

  Scenario: Getting Published area-types
    Given the following datasets based on "Example" are available
    """
    {"total_count": 1}
    """

    And I have the following population types in cantabular
    """
    ["dataset_1", "dataset_2", "dataset_3"]
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

  Scenario: Getting Unpublished area-types

    I have the following population types in cantabular
    """
    ["dataset_1", "dataset_2", "dataset_3"]
    """

    And a geography query response is available from Cantabular api extension

    And the following datasets based on "Example" are available
    """
    {"total_count": 0}
    """

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