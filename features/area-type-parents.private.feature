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
                                                },
                                                 "meta": {
                                                   "ONS_Variable": {
                                                      "Geography_Hierarchy_Order": "100"
                                                    }
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
            "items": [
                {
                    "id": "country",
                    "label": "Country",
                    "description": "",
                    "total_count": 2,
                    "hierarchy_order": 100
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

    Scenario: Getting parent area count successfully

        Given I am identified as "user@ons.gov.uk"

        And I am authorised

        Given the following parents areas count response is available from Cantabular:
        """
        {
             "Dimension": {
                "count": 1,
                "categories": [
                    {
                        "code": "E12000001",
                        "label": "Hartlepool"
                    }
                ],
                "variable": {
                    "name":  "LADCD",
                    "label": "Local Authority code"
                }
            }
        }
        """

        When I GET "/population-types/Example/area-types/city/parents/LADCD/areas-count?areas=E12000001,E12000002"

        Then I should receive the following JSON response:
        """
        1
        """

        And the HTTP status code should be "200"

    Scenario:Getting parent area count but Cantabular returns an error
        Given I am identified as "user@ons.gov.uk"

        And I am authorised

        And cantabular is unresponsive

        When I GET "/population-types/Example/area-types/city/parents/LADCD/areas-count?areas=E12000001,E12000002"

        Then I should receive the following JSON response:
        """
        {
            "errors": ["failed to get parent areas count"]
        }
        """

        And the HTTP status code should be "404"
