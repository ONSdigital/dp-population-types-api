Feature: Get Categorisations

Background:

    Given private endpoints are not enabled

    And the following categorisations response is available from Cantabular:
    """
    {
            "count": 0,
            "total_count": 0,
            "dataset": {
                "variables": {  
                    "edges": [
                    {
                        "node": {
                        "isSourceOf": {
                            "totalCount": 1,
                            "edges": [
                            {
                                "node": {
                                  "categories": {
                                   "edges": [
                                      {
                                         "node": {
                                           "meta": {
                                             "ONS_Variable": {
                                               "quality_statement_text": "quality statement"
                                               }
                                           },
                                           "code": "2",
                                           "label": "category 2"
                                          }
                                       }
                                    ]
                                    },
                                "meta": {
                                  "ONS_Variable": {
                                    "quality_statement_text": "quality statement"
                                  }
                                },
                                "label": "label 2",
                                "name": "name 2"
                                }
                            }
                            ]
                        },
                        
                        "mapFrom": []
                        }
                    }
                    ]
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

    When I GET "/population-types/Example/dimensions/hh_size/categorisations"

    Then the HTTP status code should be "200"

    And I should receive the following JSON response:
    """
        {
            "limit": 20,
	        "offset": 0,
	        "count": 0,
	        "total_count": 1,
            "items": [
                {
                    "id": "name 2",
                    "label": "label 2",
                                        "quality_statement_text": "quality statement",
                    "categories": [{
                         "id": "2",
                                             "quality_statement_text": "quality statement",
                         "label": "category 2"
                    }]
                }
            ]
        }
    """

  Scenario: Getting unpublished categorisations
    Given the following datasets based on "Example" are available
    """
    {
      "total_count": 0
    }
    """

    When I GET "/population-types/Example/dimensions/hh_size/categorisations"

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
