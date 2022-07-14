Feature: Get Area Type Parents

Background:

    Given private endpoints are not enabled

    Scenario: Getting area type parents successfully
        Given the following parents response is available from Cantabular:
        """
        {
            "dataset": {
                "variables": {
                    "edges": [
                        {
                            "node": {
                                "name":  "city",
                                "label": "City",
                                "isDirectSourceOf": {
                                    "edges": [
                                        {
                                            "node": {
                                                "name":  "country",
                                                "label": "Country",
                                                "categories": {
                                                    "totalCount": 2
                                                }
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

        When I GET "/population-types/Example/area-types/city/parents"

        Then I should receive the following JSON response:
        """
        {
            "area_types": [
                {
                    "id": "country",
                    "label": "Country",
                    "total_count": 2
                }
            ]
        }
        """

        And the HTTP status code should be "200"

    Scenario: Getting area type parents but Cantabular returns an error
        Given cantabular is unresponsive

        When I GET "/population-types/Example/area-types/city/parents"

        Then I should receive the following JSON response:
        """
        {
            "errors": ["failed to get parents: test error response"]
        }
        """

        And the HTTP status code should be "404"
