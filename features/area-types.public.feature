Feature: Area Types

  Background:
    Given private endpoints are not enabled

    And the following geography response is available from Cantabular:
    """
    {
        "count": 2,
        "total_count": 2,
        "dataset": {
          "variables": {
            "total_count": 2,
              "edges": [
                {
                  "node": {
                    "label": "Region",
                    "name":  "region",
                    "description":  "test",
                    "categories": {
                      "totalCount": 348
                    }
                  }
                },
                {
                  "node": {
                    "label": "City",
                    "name":  "city",
                    "description":  "test",
                    "categories": {
                      "totalCount": 2
                    }
                  }
                },
                {
                  "node": {
                    "name": "country",
                    "label": "Country",
                    "description": "test",
                    "categories": {
                      "totalCount": 3
                    },
                    "mapFrom":[
                      {
                        "edges":[
                          {
                            "node":{
                              "name": "city",
                              "label": "City",
                              "description": "test"
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
        "limit": 20,
        "offset": 0,
        "count": 2,
        "total_count": 2,
        "items":[
          {
                "id":"city",
                "label":"City",
                "description": "test",
                "total_count": 2
          },
          {
                "id":"country",
                "label":"Country",
                "description": "test",
                "total_count": 3
          },
          {
                "id":"region",
                "label":"Region",
                "description": "test",
                "total_count": 348
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
