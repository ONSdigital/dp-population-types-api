Feature: Population types endpoint
  As an API user
  
  I want to know all the population types for Census 2021
  
  So that I can use them to query further data

  Background:

  Given private endpoints are enabled

  And I am identified as "user@ons.gov.uk"

  And I am authorised

  And I have the following list datasets response available in cantabular
  """
  {
    "data": {
      "datasets": [
        {"name": "dataset_1"},
        {"name": "dataset_2"}
      ]
    }
  }
  """

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

  And the following datasets based on "dataset_2" are available
  """
  {
    "items": [
    {
      "id": "cantabular-flexible-default",
      "current": {
        "contacts": [
          {}
        ],
        "description": "Default Cantabular Flexible Published Collection",
        "id": "cantabular-flexible-default",
        "links": {
          "editions": {
            "href": "http://localhost:22000/datasets/cantabular-flexible-default/editions"
          },
          "latest_version": {
            "href": "http://localhost:22000/datasets/cantabular-flexible-default/editions/2021/versions/1",
            "id": "1"
          },
          "self": {
            "href": "http://localhost:22000/datasets/cantabular-flexible-default"
          }
        },
        "qmi": {},
        "state": "published",
        "title": "Cantabular Flexible Default",
        "type": "cantabular_flexible_table",
        "is_based_on": {
          "@type": "cantabular_flexible_table",
          "@id": "dataset_2"
        }
      },
      "next": {
        "contacts": [
          {}
        ],
        "description": "Default Cantabular Flexible Published Collection",
        "id": "cantabular-flexible-default",
        "links": {
          "editions": {
            "href": "http://localhost:22000/datasets/cantabular-flexible-default/editions"
          },
          "latest_version": {
            "href": "http://localhost:22000/datasets/cantabular-flexible-default/editions/2021/versions/1",
            "id": "1"
          },
          "self": {
            "href": "http://localhost:22000/datasets/cantabular-flexible-default"
          }
        },
        "qmi": {},
        "state": "published",
        "title": "Cantabular Flexible Default",
        "type": "cantabular_flexible_table",
        "is_based_on": {
          "@type": "cantabular_flexible_table",
          "@id": "dataset_2"
        }
      }
    }]
  } 
  """
  
  Scenario: The root population-types endpoint should return a list of population types
    When I GET "/population-types"

    Then I should receive the following JSON response:
    """
    {"items":[{"name": "dataset_1"}, {"name": "dataset_2"}]}
    """

  Scenario: If the root population-types endpoint fails, it should return correct errors
    Given cantabular is unresponsive
    
    When I GET "/population-types"
    
    Then I should receive the following JSON response:
    """
    {"errors": ["failed to fetch population types: cantabular failed to respond"]}
    """
