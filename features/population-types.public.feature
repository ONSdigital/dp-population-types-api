Feature: Population types endpoint
  As an API user

  I want to know all the population types for Census 2021

  So that I can use them to query further data

  Background:

  Given private endpoints are not enabled

  And I have the following population types in cantabular
  """
  {
    "datasets":[
      {
        "name": "dataset_1",
        "label": "dataset 1",
        "description": "desc_1",
        "type": "microdata"
      },
      {
        "name": "dataset_2",
        "label": "dataset 2",
        "description": "desc_2",
        "type": "tabular"
      },
      {
        "name": "dataset_3",
        "label": "dataset 3",
        "description": "desc_3",
        "type": "tabular"
      }
    ]
  }
  """

  And the following datasets based on "dataset_1" are available
  """
  {
    "total_count": 1,
    "items": [
    {
      "id": "cantabular-flexible-example"
    }]
  }
  """

  And the following datasets based on "dataset_2" are available
  """
  {
    "total_count": 1,
    "items": [
    {
      "id": "cantabular-flexible-default"
    }]
  }
  """

  Scenario: The root population-types endpoint should return a list of population types
    When I GET "/population-types"

    Then I should receive the following JSON response:
    """
    {
        "limit": 20,
        "count": 2,
        "total_count": 2,
        "offset": 0,
        "items":[
          {
            "name": "dataset_1",
            "description": "desc_1",
            "label": "dataset 1",
            "type": "microdata"
          },
          {
            "name": "dataset_2",
            "description": "desc_2",
            "label": "dataset 2",
            "type": "tabular"
          }
        ]
    }
    """

  Scenario: If the root population-types endpoint fails, it should return correct errors
    Given cantabular is unresponsive

    When I GET "/population-types"

    Then I should receive the following JSON response:
    """
    {"errors": ["failed to get population types"]}
    """

  Scenario: The population-types/{population-type} endpoint should return a population type
    When I GET "/population-types/dataset_1"

    Then I should receive the following JSON response:
    """
    {
        "population_type":{
          "name": "dataset_1",
          "description": "desc_1",
          "label": "dataset 1",
          "type": "microdata"
        }
    }
    """

    And the HTTP status code should be "200"

  Scenario: Population type not published
    When I GET "/population-types/dataset_3"

    Then I should receive the following JSON response:
    """
    {
        "errors":[
          "population type not found"
        ]
    }
    """

    And the HTTP status code should be "404"