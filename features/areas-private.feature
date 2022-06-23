Feature: Areas

  Background:
    Given private endpoints are enabled
    Given cantabular server is healthy
    And cantabular api extension is healthy
    When the following area query response is available from Cantabular api extension for the dataset "Example":
    """
    {
      "data": {
        "dataset": {
          "ruleBase": {
            "isSourceOf": {
              "search": {
                "edges": [
                  {
                    "node": {
                      "label": "City",
                      "name": "city",
                      "categories": {
                        "search": {
                          "edges": [
                            {
                              "node": {
                                "code": "0",
                                "label": "London"
                              }
                            },
                            {
                              "node": {
                                "code": "1",
                                "label": "Liverpool"
                              }
                            },
                            {
                              "node": {
                                "code": "2",
                                "label": "Belfast"
                              }
                            }
                          ]
                        }
                      }
                    }
                  }
                ]
              }
            }
          }
        }
      }
    }
    """
  Scenario: Getting areas 
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    And the following datasets based on "City" are available
    """
    {
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
            "@id": "dataset_1"
          }
        }
      }],
      "count":1
    }
    """

    And I GET "/population-types/Example/area-types/City/areas"
    Then the HTTP status code should be "200"

  Scenario: Getting areas with none published
    Given I am identified as "user@ons.gov.uk"
    And I am authorised
    And the following datasets based on "City" are available
    """
    {
      "items": [
      {
        "id": "cantabular-flexible-example",
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
          "state": "Published",
          "title": "sdf",
          "type": "cantabular_flexible_table",
          "is_based_on": {
            "@type": "cantabular_flexible_table",
            "@id": "dataset_1"
          }
        }
      }],
      "count":1
    }
    """
    And I GET "/population-types/Example/area-types/City/areas"
    Then I should receive the following JSON response:
      """
      {
        "errors": ["areas not found for private call"]
      }
      """
    And the HTTP status code should be "404"

