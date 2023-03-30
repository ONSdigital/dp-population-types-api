Feature: Area Types

  Background:
    Given private endpoints are not enabled
    And the following dimension categories response is available from Cantabular:
    """
        {
          "dataset": {
            "variables": {
              "edges": [
                {
                  "node": {
                    "categories": {
                      "edges": [
                        {
                          "node": {
                            "code": "1",
                            "label": "Female"
                          }
                        },
                        {
                          "node": {
                            "code": "2",
                            "label": "Male"
                          }
                        }
                      ],
                      "totalCount": 2
                    },
                    "label": "Sex (2 categories)",
                    "name": "sex",
                    "meta": {
                      "ONS_Variable": {
                        "quality_statement_text": "quality statement"
                      }
                    }
                  }
                }
              ]
            }
          }
        }
    """
  Scenario: Getting dimension categories happy

    When I GET "/population-types/UR/dimension-categories?dims=sex"

    Then the HTTP status code should be "200"

    And I should receive the following JSON response:
    """
    {
      "limit": 30,
      "offset": 0,
      "count": 1,
      "total_count": 1,
      "items": [
          {
              "id": "sex",
              "label": "Sex (2 categories)",
              "quality_statement_text": "quality statement",
              "default_categorisation": false,
              "categories": [
                  {
                      "id": "1",
                      "label": "Female"
                  },
                  {
                      "id": "2",
                      "label": "Male"
                  }
              ]
          }
      ]
    }
    """
  Scenario: Getting dimension categories failing
    Given the cantabular area response is bad request
    When I GET "/population-types/UR/dimension-categories?dims=sex"
    Then the HTTP status code should be "400"
    And I should receive the following JSON response:
    """
    {"errors": ["failed to get dimension categories"]}
    """
  Scenario: Getting dimension categories dataset not found
    Given the cantabular area response is not found
    When I GET "/population-types/UR/dimension-categories?dims=sex"
    Then the HTTP status code should be "404"
    And I should receive the following JSON response:
    """
    {"errors": ["failed to get dimension categories"]}
    """
