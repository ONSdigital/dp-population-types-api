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

  Scenario: 
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
