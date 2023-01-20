Feature: Dimensionas

  Background:
    Given private endpoints are enabled

    And the following dimensions response is available from Cantabular:
    """
    {
        "count":       2,
        "total_count": 2,
        "dataset": {
            "variables": {
                "search": {
                    "edges": [
                        {
                            "node": {
                                "label": "Household size (31 categories)",
                                "name":  "hh_size",
                                "description": "description",
                                "meta": {
                                    "ONS_Variable": {
                                        "quality_statement_text": "quality statement"
                                    }
                                },
                                "categories": {
                                    "totalCount": 31
                                }
                            }
                        },
                        {
                            "node": {
                                "name": "hh_tenure",
                                "label": "Tenure of household (11 categories)",
                                "description": "description",
                                "meta": {
                                    "ONS_Variable": {
                                        "quality_statement_text": "quality statement"
                                    }
                                 },
                                "categories": {
                                    "totalCount": 11
                                }
                            }
                        }
                    ]
                }
            }
        } 
    }
    """

    And the following dimensions description response is available from Cantabular:
    """
    {
        "dataset": {
            "variables": {
                "edges": [
                    {
                        "node": {
                            "categories": {
                                "totalCount": 32
                            },
                            "description": "description",
                            "meta": {
                                "ONS_Variable": {
                                    "quality_statement_text": "quality statement"
                                }
                            },
                            "label": "Number of unpaid carers in household (32 categories)",
                            "name": "hh_carers"
                        }
                    },
                    {
                        "node": {
                            "categories": {
                                "totalCount": 6
                            },
                            "description": "description",
                            "meta": {
                                "ONS_Variable": {
                                    "quality_statement_text": "quality statement"
                                }
                            },
                            "label": "Household deprivation (6 categories)",
                            "name": "hh_deprivation"
                        }
                    }
                ],
                "totalCount": 2
            }
        }
    }
    """

  Scenario: Getting published dimensions
    Given I am identified as "user@ons.gov.uk"

    And I am authorised

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

    And the following datasets based on "Example" are available
    """
    {
      "total_count": 1
    }
    """

    When I GET "/population-types/Example/dimensions"

    Then the HTTP status code should be "200"

    And I should receive the following JSON response:
    """
    {
        "limit":  20,
        "offset": 0,
        "count":  2,
        "total_count": 2,
        "items": [
            {
                "id": "hh_size",
                "label": "Household size (31 categories)",
                "description": "description",
                "quality_statement_text": "quality statement",
                "total_count": 31
            },
            {
                "id": "hh_tenure",
                "label": "Tenure of household (11 categories)",
                "quality_statement_text": "quality statement",
                "description": "description",
                "total_count": 11
            }
        ]
    }
    """

  Scenario: Getting unpublished dimensions
    Given I am identified as "user@ons.gov.uk"

    And I am authorised
    
    And I have the following population types in cantabular
    """
    {
      "datasets":[
        {
          "name": "Example1",
          "label": "Example 1"
        },
        {
          "name": "Example2",
          "label": "Example 2"
        }
      ]
    }
    """

    Given the following datasets based on "Example" are available
    """
    {
      "total_count": 1
    }
    """

    When I GET "/population-types/Example/dimensions"

    Then the HTTP status code should be "200"

    And I should receive the following JSON response:
    """
    {
        "limit":  20,
        "offset": 0,
        "count":  2,
        "total_count": 2,
        "items": [
            {
                "id": "hh_size",
                "label": "Household size (31 categories)",
                "description": "description",
                "quality_statement_text": "quality statement",
                "total_count": 31
            },
            {
                "id": "hh_tenure",
                "label": "Tenure of household (11 categories)",
                "description": "description",
                "quality_statement_text": "quality statement",
                "total_count": 11
            }
        ]
    }
    """

 Scenario: Getting dimensions description
    Given I am identified as "user@ons.gov.uk"

    And I am authorised

    And I have the following population types in cantabular
    """
    {
      "datasets":[
        {
          "name": "Example1",
          "label": "Example 1"
        },
        {
          "name": "Example2",
          "label": "Example 2"
        }
      ]
    }
    """

    And the following datasets based on "Example" are available
    """
    {
      "total_count": 1
    }
    """

    When I GET "/population-types/Example/dimensions-description?q=hh_carers&q=hh_deprivation"

    Then the HTTP status code should be "200"

    And I should receive the following JSON response:
    """
    {
        "limit":  20,
        "offset": 0,
        "count":  0,
        "total_count": 0,
        "items": [
            {
                "id": "hh_carers",
                "label": "Number of unpaid carers in household (32 categories)",
                "description": "description",
                "quality_statement_text": "quality statement",
                "total_count": 32
            },
            {
                "id": "hh_deprivation",
                "label": "Household deprivation (6 categories)",
                "description": "description",
                "quality_statement_text": "quality statement",
                "total_count": 6
            }
        ]
    }
    """
    