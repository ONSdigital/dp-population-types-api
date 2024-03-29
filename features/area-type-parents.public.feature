Feature: Get Area Type Parents

Background:

    Given private endpoints are not enabled

    And the following datasets based on "Example" are available
    """
    {
      "total_count": 1,
      "items": [
      {
        "id": "cantabular-flexible-example",
        "current": {
          "contacts": [
            {}
          ],
          "id": "cantabular-flexible-example",
          "links": {
            "editions": {
              "href": "http://localhost:22000/datasets/cantabular-flexible-example/editions"
            },
            "latest_version": {
              "href": "http://localhost:22000/datasets/cantabular-flexible-example/editions/2021/versions/1",
              "id": "1"
            },
            "self": {
              "href": "http://localhost:22000/datasets/cantabular-flexible-example"
            }
          },
          "qmi": {},
          "state": "published",
          "title": "sdf",
          "type": "cantabular_flexible_table",
          "is_based_on": {
            "@type": "cantabular_flexible_table",
            "@id": "dataset_1"
          }
        },
        "next": {
          "contacts": [
            {}
          ],
          "id": "cantabular-flexible-example",
          "links": {
            "editions": {
              "href": "http://localhost:22000/datasets/cantabular-flexible-example/editions"
            },
            "latest_version": {
              "href": "http://localhost:22000/datasets/cantabular-flexible-example/editions/2021/versions/1",
              "id": "1"
            },
            "self": {
              "href": "http://localhost:22000/datasets/cantabular-flexible-example"
            }
          },
          "qmi": {},
          "state": "published",
          "title": "sdf",
          "type": "cantabular_flexible_table",
          "is_based_on": {
            "@type": "cantabular_flexible_table",
            "@id": "Example"
          }
        }
      }]
    }
    """

    Scenario: Getting area type parents successfully

        Given the following parents response is available from Cantabular:
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
                                "description":  "test",
                                "isSourceOf": {
                                    "edges": [
                                        {
                                            "node": {
                                                "name":  "country",
                                                "label": "Country",
                                                "description":  "test",
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
            "limit": 30,
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

    Scenario: Getting area type parents on unpublished population type

        Given the following parents response is available from Cantabular:
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
                                "description":  "test",
                                "isSourceOf": {
                                    "edges": [
                                        {
                                            "node": {
                                                "name":  "country",
                                                "label": "Country",
                                                "description":  "test",
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

        And the following datasets based on "NotPublished" are available
        """
        {
          "total_count": 0
        }
        """

        When I GET "/population-types/NotPublished/area-types/city/parents"

        Then I should receive the following JSON response:
        """
        {
            "errors": ["population type not found"]
        }
        """

        And the HTTP status code should be "404"

    Scenario: Getting area type parents but Cantabular returns an error

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

        When I GET "/population-types/Example/area-types/city/parents/LADCD/areas-count?areas=E12000001,E12000002&svar=hh_tenure_9a"

        Then I should receive the following JSON response:
        """
        1
        """

        And the HTTP status code should be "200"

    Scenario:Getting parent area count but Cantabular returns an error
        Given cantabular is unresponsive

        When I GET "/population-types/Example/area-types/city/parents/LADCD/areas-count?areas=E12000001,E12000002"

        Then I should receive the following JSON response:
        """
        {
            "errors": ["failed to get parent areas count"]
        }
        """

        And the HTTP status code should be "404"
