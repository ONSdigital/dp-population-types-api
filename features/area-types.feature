Feature: Area Types

  Background:
    Given private endpoints are not enabled

  Scenario: Getting published area types

    Given the following datasets based on "Example" are available
    """
    {
      "total_count": 1
    }
    """

    And a geography query response is available from Cantabular api extension

    When I GET "/population-types/Example/area-types"

    Then the HTTP status code should be "200"

    And I should receive the following JSON response:
    """
    {
        "area-types":[
          {
                "id":"country",
                "label":"Country",
                "total_count": 2
          },
          {
                "id":"city",
                "label":"City",
                "total_count": 3
          }
        ]
    }
    """

  Scenario: Getting unpublished area types
  Given the following datasets based on "Example" are available
    """
    {
      "items": [
      {
        "id": "cantabular-flexible-example",
        "next": {
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

    And a geography query response is available from Cantabular api extension

    When I GET "/population-types/Example/area-types"

    Then the HTTP status code should be "404"

  Scenario: Getting area-types not found
    When an error json response is returned from Cantabular api extension

    And I GET "/population-types/Inexistent/area-types"

    Then I should receive the following JSON response:
    """
    {"errors":["failed to get area-types: error(s) returned by graphQL query"]}
    """

    And the HTTP status code should be "404"

  Scenario: Dataset Client returns errors
    Given the dp-dataset-api is returning errors

    And a geography query response is available from Cantabular api extension

    When I GET "/population-types/Example/area-types"

    Then the HTTP status code should be "500"

    And I should receive the following JSON response:
    """
    {"errors":["failed to get area types: failed to get population type"]}
    """
