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
        "label": "dataset 1"
      },
      {
        "name": "dataset_2",
        "label": "dataset 2"
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
            "description": "",
            "label": "dataset 1"
          },
          {
            "name": "dataset_2",
            "description": "",
            "label": "dataset 2"
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
