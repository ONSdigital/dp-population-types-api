Feature: Dimensionas

  Background:
    Given private endpoints are not enabled

    And the following dimensions response is available from Cantabular:
    """
    {
        "count":       2,
        "total_count": 2,
        "dataset": {
            "variables": {
                "search": {
                    "edges": [
                        {
                            "node": {
                                "label": "Household size (31 categories)",
                                "name":  "hh_size",
                                "description": "description",
                                "meta": {
                                    "ONS_Variable": {
                                        "quality_statement_text": "quality statement"
                                    }
                                },
                                "categories": {
                                    "totalCount": 31
                                }
                            }
                        },
                        {
                            "node": {
                                "name": "hh_tenure",
                                "label": "Tenure of household (11 categories)",
                                "description": "description",
                                "meta": {
                                    "ONS_Variable": {
                                        "quality_statement_text": "quality statement"
                                    }
                                },
                                "categories": {
                                    "totalCount": 11
                                }
                            }
                        }
                    ]
                }
            }
        } 
    }
    """

  Scenario: Getting published dimensions
    Given the following datasets based on "Example" are available
    """
    {
      "total_count": 1
    }
    """

    When I GET "/population-types/Example/dimensions"

    Then the HTTP status code should be "200"

    And I should receive the following JSON response:
    """
    {
        "limit":  30,
        "offset": 0,
        "count":  2,
        "total_count": 2,
        "items": [
            {
                "id": "hh_size",
                "label": "Household size (31 categories)",
                "description": "description",
                "quality_statement_text": "quality statement",
                "total_count": 31
            },
            {
                "id": "hh_tenure",
                "label": "Tenure of household (11 categories)",
                "description": "description",
                "quality_statement_text": "quality statement",
                "total_count": 11
            }
        ]
    }
    """

  Scenario: Getting unpublished dimensions
    Given the following datasets based on "Example" are available
    """
    {
      "total_count": 0
    }
    """

    When I GET "/population-types/Example/dimensions"

    Then I should receive the following JSON response:
    """
    {"errors":["population type not found"]}
    """

    And the HTTP status code should be "404"

  Scenario: Dataset Client returns errors
    Given the dp-dataset-api is returning errors for datasets based on "Example"

    When I GET "/population-types/Example/dimensions"

    And I should receive the following JSON response:
    """
    {"errors":["population type not found"]}
    """

    Then the HTTP status code should be "404"
