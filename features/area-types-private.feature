Feature: Area Types
  Given I am authenticated
  And these documents are in the store:
        """
        """
  And the following JSON query response is available from the Cantabular api extension:
        """
        """

  Scenario: Getting Published area-types

    When I GET "/population-types/unpublished"
    Then the service responds with status code "200" and the following JSON response:
    """
    """

  Scenario: Getting Unpublished area-types
    When I GET "/population-types/published/area-types"
    Then the service responds with status code "200" and the following JSON response:
    """
    """
