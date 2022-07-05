Feature: Area Types
  Background:
    Given private endpoints are enabled

    And the following datasets based on "dataset_1" are available
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
      }]
    }
    """
  Scenario: Getting Published area-types
    Given I am identified as "user@ons.gov.uk"

    And I am authorised

    And I have the following population types in cantabular
    """
    ["dataset_1", "dataset_2"]
    """

    And a geography query response is available from Cantabular api extension

    When I GET "/population-types/dataset_1/area-types"

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

  Scenario: Getting Unpublished area-types
   
    Given I am identified as "user@ons.gov.uk"

    And I am authorised
    
    And I have the following population types in cantabular
    """
    ["dataset_1", "dataset_2"]
    """

    And a geography query response is available from Cantabular api extension

    When I GET "/population-types/dataset_1/area-types"

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
