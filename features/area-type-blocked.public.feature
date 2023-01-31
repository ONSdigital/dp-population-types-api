Feature: Get Area Type Blocked

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

    Scenario: Getting blocked area type successfully

    Given the following blocked area response is available from Cantabular:
    """
    {
        "dataset": {
            "table": {
                "error": null,
                "rules":{ 
                    "blocked": {
                        "count": 188880
                    },
                    "evaluated": {
                        "count": 188880
                    },
                    "passed": {
                    "   count": 0
                    }
                }
            }
        }
    }
    """

    When I GET "/population-types/Example/blocked-areas-count?vars=oa,resident_age_101a"

    Then I should receive the following JSON response:
    """
    {
        "passed": 0,
        "blocked": 188880,
        "total": 188880
    }
    """

    And the HTTP status code should be "200"
    