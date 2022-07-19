Feature: Population types endpoint
  As an API user
  
  I want to know all the population types for Census 2021
  
  So that I can use them to query further data

  Background:

  Given private endpoints are not enabled

  And I have the following population types in cantabular
  """
  ["dataset_1", "dataset_2"]
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
    {"items":[{"name": "dataset_1"}, {"name": "dataset_2"}]}
    """

  Scenario: If the root population-types endpoint fails, it should return correct errors
    Given cantabular is unresponsive
    
    When I GET "/population-types"
    
    Then I should receive the following JSON response:
    """
    {"errors": ["failed to fetch population types: cantabular failed to respond"]}
    """
