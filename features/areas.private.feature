Feature: Areas

  Background:
    Given private endpoints are enabled

    Given cantabular server is healthy

    And cantabular api extension is healthy

    When the following area query response is available from Cantabular:
    """
    {
        "dataset": {
          "ruleBase": {
            "isSourceOf": {
              "search": {
                "edges": [
                  {
                    "node": {
                      "label": "City",
                      "name": "city",
                      "categories": {
                        "search": {
                          "edges": [
                            {
                              "node": {
                                "code": "0",
                                "label": "London"
                              }
                            },
                            {
                              "node": {
                                "code": "1",
                                "label": "Liverpool"
                              }
                            },
                            {
                              "node": {
                                "code": "2",
                                "label": "Belfast"
                              }
                            }
                          ]
                        }
                      }
                    }
                  }
                ]
              }
            }
          }
        }
    }
    """

  Scenario: Getting areas 
    Given I am identified as "user@ons.gov.uk"

    And I am authorised

    And I GET "/population-types/Example/area-types/City/areas"

    Then the HTTP status code should be "200"
