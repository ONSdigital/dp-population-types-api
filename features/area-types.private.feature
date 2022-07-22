Feature: Area Types
  Background:
    Given private endpoints are enabled

    And the following geography response is available from Cantabular:
    """
    {
        "dataset": {
          "ruleBase": {
            "isSourceOf": {
              "edges": [
                {
                  "node": {
                    "label": "City",
                    "name":  "city",
                    "categories": {
                      "totalCount": 3
                    }
                  }
                },
                {
                  "node": {
                    "name": "country",
                    "label": "Country",
                    "categories": {
                      "totalCount": 2
                    },
                    "mapFrom":[
                      {
                        "edges":[
                          {
                            "node":{
                              "name": "city",
<<<<<<< HEAD:features/area-types.private.feature
                              "label": "City"
=======
                              "label": "City",
>>>>>>> work for pagination:features/area-types-private.feature
                            }
                          }
                        ]
                      }
                    ]
                  }
                }
              ]
            }
          }
        }

    }
    """

  Scenario: Getting Published area-types
    Given I am identified as "user@ons.gov.uk"

    And I am authorised

    And I have the following population types in cantabular
    """
    ["dataset_1", "dataset_2"]
    """

    When I GET "/population-types/dataset_1/area-types"

    Then the HTTP status code should be "200"

    And I should receive the following JSON response:
    """
    {
        "area_types":[
          {
                "id":"city",
                "label":"City",
                "total_count": 3
          },{
                "id":"country",
                "label":"Country",
                "total_count": 2
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

    When I GET "/population-types/dataset_1/area-types"

    Then the HTTP status code should be "200"

    And I should receive the following JSON response:
    """
    {
        "area_types":[
          {
                "id":"city",
                "label":"City",
                "total_count": 3
          },{
                "id":"country",
                "label":"Country",
                "total_count": 2
          }
        ]
    }
    """
