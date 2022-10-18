Feature: Base Variables
  Background:
    Given private endpoints are enabled

    And cantabular server is healthy

    And cantabular api extension is healthy

  Scenario: Getting Base Variable Happy
    Given I am identified as "user@ons.gov.uk"

    And I am authorised

    When the following base variable response is available from Cantabular:
    """
    {
      "dataset": {
        "variables": {
          "edges": [
            {
              "node": {
                "mapFrom": [
                  {
                    "edges": [
                      {
                        "node": {
                          "label": "Accommodation type (8 categories)",
                          "name": "accommodation_type"
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
    And I GET "/population-types/dummy_data_households/dimensions/accommodation_types_5a/base"
    Then I should receive the following JSON response:
    """
    {
        "name": "accommodation_type",
        "label": "Accommodation type (8 categories)"
    }
    """
    And the HTTP status code should be "200"
