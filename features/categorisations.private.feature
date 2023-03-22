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
                                "edges": [
                                    {
                                        "node": {
                                            "categories": {
                                                "edges": [
                                                    {
                                                        "node": {
                                                            "code": "code 1",
                                                            "label": "label 1"
                                                        }
                                                    }
                                                ]
                                            },
                                            "label": "label 2",
                                            "meta": {
                                                "Default_Classification_Flag": "Y",
                                                "ONS_Variable": {
                                                    "Quality_Statement_Text": "quality statement 1"
                                                }
                                            },
                                            "name": "name 1"
                                        }
                                    }
                                ],
                                "totalCount": 1
                            }
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
                    "id": "name 1",
                    "label": "label 2",
                    "quality_statement_text":"quality statement 1",
                    "default_categorisation": true,
                    "categories": [{
                         "id": "code 1",
                         "label": "label 1"
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