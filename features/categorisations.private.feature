Feature: Get Categorisations

Background:

    Given private endpoints are enabled

    Scenario: Getting categorisations successfully

        Given I am identified as "user@ons.gov.uk"

        And I am authorised

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
                                           "code": "1",
                                           "label": "category 1"
                                          }
                                       }
                                    ]
                                    },
                                "meta": {
                                  "ONS_Variable": {
                                    "quality_statement_text": "quality statement"
                                  }
                                },
                                "label": "label 1",
                                "name": "code 1"
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

        When I GET "/population-types/Example/dimensions/hh_size/categorisations"

        Then I should receive the following JSON response:
        """
        {
            "limit": 20,
	        "offset": 0,
	        "count": 0,
	        "total_count": 1,
            "items": [
                {
                    "id": "code 1",
                    "label": "label 1",
                    "quality_statement_text": "quality statement",
                    "categories": [{
                         "id": "1",
                         "label": "category 1",
                         "quality_statement_text": "quality statement"
                    }]
                }
            ]
        }
        """

        And the HTTP status code should be "200"

    Scenario: Getting categorisations but Cantabular returns an error
        Given I am identified as "user@ons.gov.uk"

        And I am authorised

        And cantabular is unresponsive

        When I GET "/population-types/Example/dimensions/hh_size/categorisations"

        Then I should receive the following JSON response:
        """
        {
            "errors": ["failed to get categorisations"]
        }
        """

        And the HTTP status code should be "404"