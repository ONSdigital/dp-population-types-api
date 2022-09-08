Feature: Get Categorisations

Background:

    Given private endpoints are enabled

    Scenario: Getting categorisations successfully

        Given I am identified as "user@ons.gov.uk"

        And I am authorised

        And the following categorisations response is available from Cantabular:
        """
        {
            "count": 1,
            "total_count": 1,
            "dataset": {
                "variables": {  
                "search": {
                    "edges": [
                    {
                        "node": {
                        "categories": {
                            "edges": [
                            {
                                "node": {
                                "label": "label 1",
                                "code": "code 1"
                                }
                            }
                            ]
                        },
                        "name": "name 2",
                        "label": "label 2"
                        }
                    }
                    ]
                }
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
	        "count": 1,
	        "total_count": 1,
            "items": [
                {
                    "name": "name 2",
                    "label": "label 2"
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