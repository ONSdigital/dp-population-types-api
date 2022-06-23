Feature: Areas

  Background:
    Given private endpoints are not enabled
    And I am not authorised
    And cantabular server is healthy
    And cantabular api extension is healthy
    And the following datasets based on "City" are available
    """
    {"count": 1}
    """
  Scenario: Getting areas happy
    When the following area query response is available from Cantabular api extension for the dataset "Example":
      """
      {
        "data": {
          "dataset": {
            "ruleBase": {
              "isSourceOf": {
                "search": {
                  "edges": [
                    {
                      "node": {
                        "label": "City",
                        "name": "city",
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
                          }
                        }
                      }
                    }
                  ]
                }
              }
            }
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
            "area-type": "city"
          },
          {
            "id": "1",
            "label": "Liverpool",
            "area-type": "city"
          },
          {
            "id": "2",
            "label": "Belfast",
            "area-type": "city"
          }
        ]
      }
      """
    And the HTTP status code should be "200"

  Scenario: Getting areas specific search
    When the following area query response is available from Cantabular api extension for the dataset "Example", area type "City" and text "London":
    """
    {
      "data": {
        "dataset": {
          "ruleBase": {
            "isSourceOf": {
              "search": {
                "edges": [
                  {
                    "node": {
                      "label": "City",
                      "name": "city",
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
                        }
                      }
                    }
                  }
                ]
              }
            }
          }
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
          "area-type": "city"
        }
      ]
    }
    """

  Scenario: Getting areas no dataset or search text
    When the following area query response is available from Cantabular api extension for the dataset "", area type "" and text "":
    """
    {
      "data": {
        "dataset": {
          "ruleBase": {
            "isSourceOf": {
              "search": {
                "edges": [
                  {
                    "node": {
                      "label": "Country",
                      "name": "country",
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
                        }
                      }
                    }
                  },
                  {
                    "node": {
                      "label": "City",
                      "name": "city",
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
                        }
                      }
                    }
                  }
                ]
              }
            }
          }
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
          "area-type": "country"
        },
        {
          "id": "N",
          "label": "Northern Ireland",
          "area-type": "country"
        },
        {
          "id": "0",
          "label": "London",
          "area-type": "city"
        },
        {
          "id": "1",
          "label": "Liverpool",
          "area-type": "city"
        },
        {
          "id": "2",
          "label": "Belfast",
          "area-type": "city"
        }
      ]
    }
    """

  Scenario: Getting areas invalid dataset
    Given cantabular is unresponsive
    When the following area query response is available from Cantabular api extension for the dataset "Test":
    """
    {
      "data": {
        "dataset": null
      },
      "errors": [
        {
          "message": "404 Not Found: dataset not loaded in this server",
          "locations": [
            {
              "line": 2,
              "column": 2
            }
          ],
          "path": [
            "dataset"
          ]
        }
      ]
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
    When the following area query response is available from Cantabular api extension for the dataset "Example", area type "" and text "rio":
    """
    {
      "data": {
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
    When the following area query response is available from Cantabular api extension for the dataset "Example":
      """
      {
        "data": {
          "dataset": {
            "ruleBase": {
              "isSourceOf": {
                "search": {
                  "edges": [
                    {
                      "node": {
                        "label": "City2",
                        "name": "city",
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
                          }
                        }
                      }
                    }
                  ]
                }
              }
            }
          }
        }
      }
      """
    
    When the following datasets based on "City2" are available
    """
    {"count": 0}
    """
    And I GET "/population-types/Example/area-types/City2/areas"
    Then I should receive the following JSON response:
      """
      {
        "errors": ["areas not found"]
      }
      """
    And the HTTP status code should be "404"