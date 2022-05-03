Feature: Area Types

  Scenario: Getting area-types happy
    When a geography query response is available from Cantabular api extension
    And I GET "/population-types/Example/area-types"
    Then a list of area types is returned
    And the HTTP status code should be "200"


  Scenario: Getting area-types not found
    When an error json response is returned from Cantabular api extension
    And I GET "/population-types/Inexistent/area-types"
    Then the service responds with "404" http code and an internal server error saying "failed to get area-types: error(s) returned by graphQL query"
    And the HTTP status code should be "404"
