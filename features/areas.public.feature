Feature: Areas

  Background:
    Given private endpoints are not enabled

    And cantabular server is healthy

    And cantabular api extension is healthy

    And the following datasets based on "Example" are available
    """
    {"total_count": 1}
    """

  Scenario: Getting areas happy

    When the following area query response is available from Cantabular:
    """
    {
      "dataset": {
        "variables": {
          "edges": [
            {
              "node": {
                "categories": {
                  "search": {
                    "edges": [
                      {
                        "node": {
                            "code": "0",
                            "label": "London"
                        }
                      },
                      {
                        "node": {
                          "code": "1",
                          "label": "Liverpool"
                        }
                      },
                      {
                        "node": {
                          "code": "2",
                          "label": "Belfast"
                        }
                      }
                    ]
                  },
                  "totalCount": 3
                },
                "label": "City",
                "name": "city"
              }
            }
          ]
        }
      }
    }
    """

    And I GET "/population-types/Example/area-types/City/areas"

    Then I should receive the following JSON response:
      """
      {
        "areas": [
          {
            "id": "0",
            "label": "London",
            "area_type": "city"
          },
          {
            "id": "1",
            "label": "Liverpool",
            "area_type": "city"
          },
          {
            "id": "2",
            "label": "Belfast",
            "area_type": "city"
          }
        ]
      }
      """
    And the HTTP status code should be "200"

  Scenario: Getting areas specific search
    
    And the following area query response is available from Cantabular:
    """
    {
      "dataset": {
        "variables": {
          "edges": [
            {
              "node": {
                "categories": {
                  "search": {
                    "edges": [
                        {
                          "node": {
                            "code": "0",
                            "label": "London"
                          }
                        }
                      ]
                    },
                    "totalCount": 1
                  },
                  "label": "City",
                  "name": "city"
                }
              }
          ]
        }
      }
    }
    """

    And I GET "/population-types/Example/area-types/City/areas?q=London"

    Then I should receive the following JSON response:
    """
    {
      "areas": [
        {
          "id": "0",
          "label": "London",
          "area_type": "city"
        }
      ]
    }
    """

  Scenario: Getting areas no dataset or search text
    Given I am identified as "user@ons.gov.uk"

    And I am authorised
    
    And the following area query response is available from Cantabular:
    """
    {
      "dataset": {
        "variables": {
          "edges": [
            {
              "node": {
                "categories": {
                  "search": {
                    "edges": [
                      {
                        "node": {
                          "code": "E",
                          "label": "England"
                        }
                      },
                      {
                        "node": {
                            "code": "N",
                            "label": "Northern Ireland"
                        }
                      }
                    ]
                  },
                  "totalCount": 2
                },
                "label": "Country",
                "name": "country"
              }
            },
            {
              "node": {
                "categories": {
                  "search": {
                    "edges": [
                      {
                        "node": {
                          "code": "0",
                          "label": "London"
                        }
                      },
                      {
                        "node": {
                          "code": "1",
                          "label": "Liverpool"
                        }
                      },
                      {
                        "node": {
                          "code": "2",
                          "label": "Belfast"
                        }
                      }
                    ]
                  },
                  "totalCount": 3
                },
                "label": "City",
                "name": "city"
              }
            }
          ]
        }
      }
    }

    """

    And I GET "/population-types/Example/area-types/City/areas"

    Then I should receive the following JSON response:
    """
    {
      "areas": [
        {
          "id": "E",
          "label": "England",
          "area_type": "country"
        },
        {
          "id": "N",
          "label": "Northern Ireland",
          "area_type": "country"
        },
        {
          "id": "0",
          "label": "London",
          "area_type": "city"
        },
        {
          "id": "1",
          "label": "Liverpool",
          "area_type": "city"
        },
        {
          "id": "2",
          "label": "Belfast",
          "area_type": "city"
        }
      ]
    }
    """

  Scenario: Getting areas invalid Dataset

    And cantabular is unresponsive

    When the following area query response is available from Cantabular:
    """
    {
      "dataset": null
    }
    """
    
    And I GET "/population-types/Example/area-types/City/areas"
    
    Then I should receive the following JSON response:
    """
    {
      "errors": [
        "failed to get areas"
      ]
    }
    """

  Scenario: Get areas area does not exist
    
    And the following area query response is available from Cantabular:
    """
    {
        "dataset": {
          "ruleBase": {
            "isSourceOf": {
              "search": {
                "edges": []
              }
            }
          }
        }

    }
    """

    And I GET "/population-types/Example/area-types/City/areas?q=rio"

    Then I should receive the following JSON response:
    """
    {
      "areas": null
    }
    """

  Scenario: Getting areas when dataset is not available
    
    And the following area query response is available from Cantabular:
    """
    {
      "dataset": {
        "variables": {
          "edges": [
            {
              "node": {
                "categories": {
                  "search": {
                    "edges": [
                      {
                        "node": {
                            "code": "0",
                            "label": "London"
                        }
                      },
                      {
                        "node": {
                          "code": "1",
                          "label": "Liverpool"
                        }
                      },
                      {
                        "node": {
                          "code": "2",
                          "label": "Belfast"
                        }
                      }
                    ]
                  },
                  "totalCount": 3
                },
                "label": "City2",
                "name": "city"
              }
            }
          ]
        }
      }
    }
    """
    
    And the following datasets based on "Example2" are available
    """
    {"count": 0}
    """

    And I GET "/population-types/Example2/area-types/City/areas"

    Then I should receive the following JSON response:
      """
      {
        "errors": ["population type not found"]
      }
      """
    And the HTTP status code should be "404"