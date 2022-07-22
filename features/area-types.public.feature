Feature: Area Types

  Background:
    Given private endpoints are not enabled

    And the following geography response is available from Cantabular:
    """
    {

        "dataset": {
          "ruleBase": {
            "isSourceOf": {
              "edges": [
                {
                  "node": {
                    "label": "City",
                    "name":  "city",
                    "categories": {
                      "totalCount": 3
                    }
                  }
                },
                {
                  "node": {
                    "name": "country",
                    "label": "Country",
                    "categories": {
                      "totalCount": 2
                    },
                    "mapFrom":[
                      {
                        "edges":[
                          {
                            "node":{
                              "name": "city",
                              "label": "City",
                              "filterOnly": "false"
                            }
                          }
                        ]
                      }
                    ]
                  }
                }
              ]
            }
          }
        }

    }
    """

  Scenario: Getting published area types

    Given the following datasets based on "Example" are available
    """
    {
      "total_count": 1
    }
    """

    When I GET "/population-types/Example/area-types"

    Then the HTTP status code should be "200"

    And I should receive the following JSON response:
    """
    {
        "area_types":[
          {
                "id":"city",
                "label":"City",
                "total_count": 3
          },
          {
                "id":"country",
                "label":"Country",
                "total_count": 2
          }
        ]
    }
    """

  Scenario: Getting unpublished area types
    Given the following datasets based on "Example" are available
    """
    {
      "total_count": 0
    }
    """

    When I GET "/population-types/Example/area-types"

    Then I should receive the following JSON response:
    """
    {"errors":["population type not found"]}
    """

    And the HTTP status code should be "404"

  Scenario: Dataset Client returns errors
    Given the dp-dataset-api is returning errors for datasets based on "Example"

    When I GET "/population-types/Example/area-types"

    And I should receive the following JSON response:
    """
    {"errors":["population type not found"]}
    """

    Then the HTTP status code should be "404"
