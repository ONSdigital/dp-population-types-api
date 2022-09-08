Feature: Get Area Type Parents

Background:

    Given private endpoints are enabled

    Scenario: Getting area type parents successfully

        Given I am identified as "user@ons.gov.uk"

        And I am authorised

        And the following parents response is available from Cantabular:
        """
        {
            "count": 1,
            "total_count": 1,
            "dataset": {
                "variables": {
                    "edges": [
                        {
                            "node": {
                                "name":  "city",
                                "label": "City",
                                "isSourceOf": {
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
            "limit": 20,
            "offset": 0,
            "count": 1,
            "total_count": 1,
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
        Given I am identified as "user@ons.gov.uk"

        And I am authorised

        And cantabular is unresponsive

        When I GET "/population-types/Example/area-types/city/parents"

        Then I should receive the following JSON response:
        """
        {
            "errors": ["failed to get parents"]
        }
        """

        And the HTTP status code should be "404"
