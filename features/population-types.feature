Feature: Population types endpoint
  As an API user
  I want to know all the population types for Census 2021
  So that I can use them to query further data

  Background:
    Given private endpoints are not enabled

  Scenario: The root population-types endpoint should return a list of population types
    Given I have some population types in cantabular
    When I access the root population types endpoint
    Then a list of named cantabular population types is returned

  Scenario: If the root population-types endpoint fails, it should return correct errors
    Given cantabular is unresponsive
    When I access the root population types endpoint
    Then an internal server error saying "failed to fetch population types" is returned
    And the HTTP status code should be "500"
