Feature: Base Variables
  Background:
    Given private endpoints are not enabled

    And cantabular server is healthy

    And cantabular api extension is healthy

    And the following datasets based on "Example" are available
    """
    {"total_count": 1}
    """
  Scenario: Getting Base Variable Happy
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
    And I GET "/population-types/Example/dimensions/accommodation_type_5a/base"
    Then I should receive the following JSON response:
    """
    {
        "name": "accommodation_type",
        "label": "Accommodation type (8 categories)"
    }
    """
    And the HTTP status code should be "200"

  Scenario: Variable Not Found
    When the cantabular response is bad gateway
    And I GET "/population-types/Example/dimensions/NotExists/base"
    Then I should receive the following JSON response:
    """
    {
    "errors": ["failed to get base variable"]
    }
    """
    # NB: service returns 502 for a dimensions not found due
    # to inconsistent error messages from cantabular and the setup
    # of the cantabular api client
    # so real world behaviour was mocked here

    # Should be changed when sensible code changes error messages to be more consistent
    And the HTTP status code should be "502"

  Scenario: Dataset Not Found
    When the cantabular area response is not found
    And I GET "/population-types/NOTEXISTS/dimensions/accommodation_types_5a/base"
    Then I should receive the following JSON response:
    """
    {
    "errors": ["failed to get base variable"]
    }
    """
    And the HTTP status code should be "404"
