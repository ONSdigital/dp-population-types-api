Feature: Single Area

  Background:
    Given private endpoints are not enabled

    And cantabular server is healthy

    And cantabular api extension is healthy

    And the following datasets based on "Example" are available
    """
    {"total_count": 1}
    """

    Scenario: Getting area happy

        When the following area query response is available from Cantabular:
        """
        {
        "dataset": {
          "variables": {
            "edges": [
              {
                "node": {
                  "categories": {
                    "search": {
                      "edges": [
                        {
                          "node": {
                            "code": "1",
                            "label": "Liverpool"
                          }
                        }
                      ]
                    },
                    "totalCount": 1
                  },
                  "label": "City",
                  "name": "city"
                }
              }
            ]
          }
        }
        }
        """

        And I GET "/population-types/Example/area-types/City/areas/Liverpool"
        Then I should receive the following JSON response:
          """
          {"area":
          {
            "id": "1",
            "label": "Liverpool",
            "area_type": "city"
          }
          }
          """
        And the HTTP status code should be "200"
    Scenario: Area Not Found
        When the cantabular area response is bad request
        And I GET "/population-types/NOTEXIST/area-types/City/areas/Liverpool"
        Then the HTTP status code should be "400"
        And I should receive the following JSON response:
        """
        {
        "errors": ["failed to get area"]
        }
        """
    Scenario: Variable Not Found
        When the cantabular area response is bad request
        And I GET "/population-types/Example/area-types/NOTEXIST/areas/Liverpool"
        Then the HTTP status code should be "400"
        And I should receive the following JSON response:
        """
        {
        "errors": ["failed to get area"]
        }
        """

    Scenario: Code Not Found
        When the cantabular area response is not found
        And I GET "/population-types/Example/area-types/City/areas/NOTEXIST"
        Then the HTTP status code should be "404"
        And I should receive the following JSON response:
        """
        {
        "errors": ["failed to get area"]
        }
        """
    Scenario: Partials not matched
        When the cantabular area response is not found
        And I GET "/population-types/Example/area-types/City/areas/Be"
        Then the HTTP status code should be "404"
        And I should receive the following JSON response:
        """
        {
        "errors": ["failed to get area"]
        }
        """
