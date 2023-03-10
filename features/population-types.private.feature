Feature: Population types endpoint
  As an API user

  I want to know all the population types for Census 2021

  So that I can use them to query further data

  Background:

  Given private endpoints are enabled

  And I am identified as "user@ons.gov.uk"

  And I am authorised

  And I have the following population types in cantabular
  """
  {
    "datasets":[
      {
        "name": "dataset_1",
        "description": "dataset_1",
        "label": "dataset 1",
        "type": "microdata"
      },
      {
        "name": "dataset_2",
        "description": "dataset_2",
        "label": "dataset 2",
        "type": "tabular"
      }
    ]
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
            "description": "dataset_1",
            "label": "dataset 1",
            "type": "microdata"
          },
          {
            "name": "dataset_2",
            "description": "dataset_2",
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
  Scenario: The root population-types endpoint should return only one population type
    Given I have this metadata:
    """
    [
      {
        "_id": "dataset_1",
        "default_dataset_id": "default-dataset",
        "id": "dataset_1"
      }
    ]
    """
    When I GET "/population-types?require-default-dataset=true"
    Then I should receive the following JSON response:
    """
    {
        "limit": 20,
        "count": 1,
        "total_count": 1,
        "offset": 0,
        "items":[
          {
            "name": "dataset_1",
            "description": "dataset_1",
            "label": "dataset 1",
            "type": "microdata"
          }
        ]
    }
    """

  Scenario: The population-types/{population-type} endpoint should return a population type
    When I GET "/population-types/dataset_1"

    Then I should receive the following JSON response:
    """
    {
        "population_type":{
          "name": "dataset_1",
          "description": "dataset_1",
          "label": "dataset 1",
          "type": "microdata"
        }
    }
    """

    And the HTTP status code should be "200"

  Scenario: Population type not found
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
