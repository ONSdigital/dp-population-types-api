Feature: Area Types
  Background:
    Given private endpoints are enabled

    And the following geography response is available from Cantabular:
    """
    {
        "count": 2,
        "total_count": 2,
        "dataset": {
          "variables": {
          "total_count": 2,
            "edges": [
              {
                "node": {
                  "label": "City",
                   "name":  "city",
                   "description":  "test",
                   "categories": {
                     "totalCount": 3
                  },
                  "meta": {
                    "ONS_Variable": {
                      "Geography_Hierarchy_Order": "100"
                     }
                   }
                }
              },
              {
                "node": {
                  "name": "country",
                  "label": "Country",
                   "description":  "test",
                  "categories": {
                    "totalCount": 2
                   },
                   "mapFrom":[
                     {
                       "edges":[
                         {
                           "node":{
                             "name": "city",
                             "label": "City",
                             "description":  "test",
                             "filterOnly": "false"
                           }
                         }
                       ]
                     }
                   ],
                   "meta": {
                     "ONS_Variable": {
                       "Geography_Hierarchy_Order": "200"
                     }
                   }
                 }
               }
             ]
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
        "limit": 20,
        "offset": 0,
        "count": 2,
        "total_count": 2,
        "items":[
          {
                "id":"country",
                "label":"Country",
                "description": "test",
                "total_count": 2,
                "geography_hierarchy_order": 200
          },
          {
                "id":"city",
                "label":"City",
                "description": "test",
                "total_count": 3,
                "geography_hierarchy_order": 100
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
        "limit": 20,
        "offset": 0,
        "count": 2,
        "total_count": 2,
        "items":[
          {
                "id":"country",
                "label":"Country",
                "description": "test",
                "total_count": 2,
                "geography_hierarchy_order": 200
          },
          {
                "id":"city",
                "label":"City",
                "description": "test",
                "total_count": 3,
                "geography_hierarchy_order": 100
          }
        ]
    }
    """
