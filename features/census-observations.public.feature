Feature: Get census observations

Background:

    Given private endpoints are not enabled

Scenario: Census observations endpoint is not enabled

    When I GET "/population-types/UR/census-observations?dimensions=ltla,resident_age_7b&area-type=ltla,E06000001"

    And the HTTP status code should be "404"

Scenario: Getting census observations successfully

    Given census observations endpoint is enabled

    And the following dataset type is available from Cantabular:
    """
    {
        "data": {
            "dataset": {
            "type": "microdata"
            }
        }
    }
    """

    And the following census observations response is available from Cantabular:
    """
    {
        "dataset": {
            "table": {
                "dimensions": [
                    {
                        "categories": [
                            {
                                "code": "E06000001",
                                "label": "Hartlepool"
                            }
                        ],
                        "count": 1,
                        "variable": {
                        "label": "Lower Tier Local Authorities",
                        "name": "ltla"
                        }
                    },
                    {
                        "categories": [
                            {
                                "code": "1",
                                "label": "Aged 64 years and under"
                            },
                            {
                                "code": "2",
                                "label": "Aged 65 to 69 years"
                            },
                            {
                                "code": "3",
                                "label": "Aged 70 to 74 years"
                            },
                            {
                                "code": "4",
                                "label": "Aged 75 to 79 years"
                            },
                            {
                                "code": "5",
                                "label": "Aged 80 to 84 years"
                            },
                            {
                                "code": "6",
                                "label": "Aged 85 to 89 years"
                            },
                            {
                                "code": "7",
                                "label": "Aged 90 years and over"
                            }
                        ],
                        "count": 7,
                        "variable": {
                        "label": "Age (B) (7 categories)",
                        "name": "resident_age_7b"
                        }
                    }
                ],
                "error": null,
                "rules": {
                    "blocked": {
                        "count": 0
                    },
                    "evaluated": {
                        "count": 1
                    },
                    "passed": {
                        "count": 1
                    }
                },
                "values": [
                    57326,
                    4376,
                    4311,
                    4355,
                    4345,
                    4449,
                    22878
                ]
            }
        }
    }
    """

    And the following JSON response is available to stream:
    """
    {
    "observations": [
        {
            "dimensions": [
                {
                    "dimension": "Lower Tier Local Authorities",
                    "dimension_id": "ltla",
                    "option": "Hartlepool",
                    "option_id": "E06000001"
                },
                {
                    "dimension": "Age (B) (7 categories)",
                    "dimension_id": "resident_age_7b",
                    "option": "Aged 64 years and under",
                    "option_id": "1"
                }
            ],
            "observation": 57326
        },
        {
            "dimensions": [
                {
                    "dimension": "Lower Tier Local Authorities",
                    "dimension_id": "ltla",
                    "option": "Hartlepool",
                    "option_id": "E06000001"
                },
                {
                    "dimension": "Age (B) (7 categories)",
                    "dimension_id": "resident_age_7b",
                    "option": "Aged 65 to 69 years",
                    "option_id": "2"
                }
            ],
            "observation": 4376
        },
        {
            "dimensions": [
                {
                    "dimension": "Lower Tier Local Authorities",
                    "dimension_id": "ltla",
                    "option": "Hartlepool",
                    "option_id": "E06000001"
                },
                {
                    "dimension": "Age (B) (7 categories)",
                    "dimension_id": "resident_age_7b",
                    "option": "Aged 70 to 74 years",
                    "option_id": "3"
                }
            ],
            "observation": 4311
        },
        {
            "dimensions": [
                {
                    "dimension": "Lower Tier Local Authorities",
                    "dimension_id": "ltla",
                    "option": "Hartlepool",
                    "option_id": "E06000001"
                },
                {
                    "dimension": "Age (B) (7 categories)",
                    "dimension_id": "resident_age_7b",
                    "option": "Aged 75 to 79 years",
                    "option_id": "4"
                }
            ],
            "observation": 4355
        },
        {
            "dimensions": [
                {
                    "dimension": "Lower Tier Local Authorities",
                    "dimension_id": "ltla",
                    "option": "Hartlepool",
                    "option_id": "E06000001"
                },
                {
                    "dimension": "Age (B) (7 categories)",
                    "dimension_id": "resident_age_7b",
                    "option": "Aged 80 to 84 years",
                    "option_id": "5"
                }
            ],
            "observation": 4345
        },
        {
            "dimensions": [
                {
                    "dimension": "Lower Tier Local Authorities",
                    "dimension_id": "ltla",
                    "option": "Hartlepool",
                    "option_id": "E06000001"
                },
                {
                    "dimension": "Age (B) (7 categories)",
                    "dimension_id": "resident_age_7b",
                    "option": "Aged 85 to 89 years",
                    "option_id": "6"
                }
            ],
            "observation": 4449
        },
        {
            "dimensions": [
                {
                    "dimension": "Lower Tier Local Authorities",
                    "dimension_id": "ltla",
                    "option": "Hartlepool",
                    "option_id": "E06000001"
                },
                {
                    "dimension": "Age (B) (7 categories)",
                    "dimension_id": "resident_age_7b",
                    "option": "Aged 90 years and over",
                    "option_id": "7"
                }
            ],
            "observation": 22878
        }
    ],
    "total_observations": 7,
    "blocked_areas": 0,
    "total_areas": 1,
    "areas_returned": 1
    }
    """


    When I GET "/population-types/UR/census-observations?dimensions=ltla,resident_age_7b&area-type=ltla,E06000001"

    Then I should receive the following JSON response:
    """
    {
        "observations": [
            {
                "dimensions": [
                    {
                        "dimension": "Lower Tier Local Authorities",
                        "dimension_id": "ltla",
                        "option": "Hartlepool",
                        "option_id": "E06000001"
                    },
                    {
                        "dimension": "Age (B) (7 categories)",
                        "dimension_id": "resident_age_7b",
                        "option": "Aged 64 years and under",
                        "option_id": "1"
                    }
                ],
                "observation": 57326
            },
            {
                "dimensions": [
                    {
                        "dimension": "Lower Tier Local Authorities",
                        "dimension_id": "ltla",
                        "option": "Hartlepool",
                        "option_id": "E06000001"
                    },
                    {
                        "dimension": "Age (B) (7 categories)",
                        "dimension_id": "resident_age_7b",
                        "option": "Aged 65 to 69 years",
                        "option_id": "2"
                    }
                ],
                "observation": 4376
            },
            {
                "dimensions": [
                    {
                        "dimension": "Lower Tier Local Authorities",
                        "dimension_id": "ltla",
                        "option": "Hartlepool",
                        "option_id": "E06000001"
                    },
                    {
                        "dimension": "Age (B) (7 categories)",
                        "dimension_id": "resident_age_7b",
                        "option": "Aged 70 to 74 years",
                        "option_id": "3"
                    }
                ],
                "observation": 4311
            },
            {
                "dimensions": [
                    {
                        "dimension": "Lower Tier Local Authorities",
                        "dimension_id": "ltla",
                        "option": "Hartlepool",
                        "option_id": "E06000001"
                    },
                    {
                        "dimension": "Age (B) (7 categories)",
                        "dimension_id": "resident_age_7b",
                        "option": "Aged 75 to 79 years",
                        "option_id": "4"
                    }
                ],
                "observation": 4355
            },
            {
                "dimensions": [
                    {
                        "dimension": "Lower Tier Local Authorities",
                        "dimension_id": "ltla",
                        "option": "Hartlepool",
                        "option_id": "E06000001"
                    },
                    {
                        "dimension": "Age (B) (7 categories)",
                        "dimension_id": "resident_age_7b",
                        "option": "Aged 80 to 84 years",
                        "option_id": "5"
                    }
                ],
                "observation": 4345
            },
            {
                "dimensions": [
                    {
                        "dimension": "Lower Tier Local Authorities",
                        "dimension_id": "ltla",
                        "option": "Hartlepool",
                        "option_id": "E06000001"
                    },
                    {
                        "dimension": "Age (B) (7 categories)",
                        "dimension_id": "resident_age_7b",
                        "option": "Aged 85 to 89 years",
                        "option_id": "6"
                    }
                ],
                "observation": 4449
            },
            {
                "dimensions": [
                    {
                        "dimension": "Lower Tier Local Authorities",
                        "dimension_id": "ltla",
                        "option": "Hartlepool",
                        "option_id": "E06000001"
                    },
                    {
                        "dimension": "Age (B) (7 categories)",
                        "dimension_id": "resident_age_7b",
                        "option": "Aged 90 years and over",
                        "option_id": "7"
                    }
                ],
                "observation": 22878
            }
        ],
        "links": {
            "self": {
                "href": "http://foo/population-types/UR/census-observations?dimensions=ltla,resident_age_7b&area-type=ltla,E06000001"
            }
        },
        "total_observations": 7,
        "blocked_areas": 0,
        "total_areas": 1,
        "areas_returned": 1
    }
    """

    And the HTTP status code should be "200"

Scenario: Getting census observations error

Given the following census observations response is available from Cantabular:
    """
    {
        "dataset": {
        "table": null
        }
    }
    """
    
    And the cantabular area response is bad request

    And the following dataset type is available from Cantabular:
    """
    {
        "data": {
            "dataset": {
            "type": "microdata"
            }
        }
    }
    """

    And the following JSON response is available to stream:
    """
    {
    "errors": [
        "transform error: error(s) returned by graphQL query: 400 Bad Request: variable for filter at position 1 does not exist"
    ]
    }
    """

    And census observations endpoint is enabled

    When I GET "/population-types/UR/census-observations?dimensions=ltla,resident_age_7b&area-type=ltla,E06000001"

    Then the HTTP status code should be "400"


Scenario: Getting census observations dataset error

Given the following dataset type is available from Cantabular:
    """
    {
        "data": {
            "dataset": {
            "type": "not-microdata"
            }
        }
    }
    """

    And census observations endpoint is enabled

    When I GET "/population-types/UR/census-observations?dimensions=ltla,resident_age_7b&area-type=ltla,E06000001"

    Then the HTTP status code should be "400"

Scenario: Getting More Than 5 errors

    Given the following census observations response is available from Cantabular:
    """
    {
        "data": {
          "dataset": {
            "table": {
                "dimensions": null,
                "error": "Maximum variables at MSOA and above is 5",
                "rules": null,
                "values": null
            }
           }
        }
    }
    """

    And the following dataset type is available from Cantabular:
    """
    {
        "data": {
            "dataset": {
            "type": "microdata"
            }
        }
    }
    """

    And the count check returns a bad request
    
    And census observations endpoint is enabled

    When I GET "/population-types/UR/census-observations?dimensions=ltla,resident_age_7b&area-type=ltla,E06000001"
    Then the HTTP status code should be "400"
    Then I should receive the following JSON response:
    """
    {
    "errors": [
        "Maximum variables at MSOA and above is 5"
    ]
    }
    """

Scenario: Too many observations returned

    Given the following census observations response is available from Cantabular:
     """
    {
        "dataset": {
            "table": {
                "dimensions": [
                    {
                        "categories": [
                            {
                                "code": "E06000001",
                                "label": "Hartlepool"
                            }
                        ],
                        "count": 1,
                        "variable": {
                        "label": "Lower Tier Local Authorities",
                        "name": "ltla"
                        }
                    },
                    {
                        "categories": [
                            {
                                "code": "1",
                                "label": "Aged 64 years and under"
                            },
                            {
                                "code": "2",
                                "label": "Aged 65 to 69 years"
                            },
                            {
                                "code": "3",
                                "label": "Aged 70 to 74 years"
                            },
                            {
                                "code": "4",
                                "label": "Aged 75 to 79 years"
                            },
                            {
                                "code": "5",
                                "label": "Aged 80 to 84 years"
                            },
                            {
                                "code": "6",
                                "label": "Aged 85 to 89 years"
                            },
                            {
                                "code": "7",
                                "label": "Aged 90 years and over"
                            }
                        ],
                        "count": 7,
                        "variable": {
                        "label": "Age (B) (7 categories)",
                        "name": "resident_age_7b"
                        }
                    }
                ],
                "error": null,
                "rules": {
                    "blocked": {
                        "count": 0
                    },
                    "evaluated": {
                        "count": 1
                    },
                    "passed": {
                        "count": 1
                    }
                },
                "values": [
                    57326,
                    4376,
                    4311,
                    4355,
                    4345,
                    4449,
                    22878
                ]
            }
        }
    }
    """

    And the following dataset type is available from Cantabular:
    """
    {
        "data": {
            "dataset": {
            "type": "microdata"
            }
        }
    }
    """

    And the count check returns that the count is too large
    
    And census observations endpoint is enabled

    When I GET "/population-types/UR/census-observations?dimensions=ltla,resident_age_7b&area-type=ltla,E06000001"
    Then the HTTP status code should be "403"
    Then I should receive the following JSON response:
    """
    {
    "errors": [
        "Too many rows returned, please refine your query by requesting specific areas or reducing the number of categories returned.  For further information please visit https://developer.ons.gov.uk/createyourowndataset/"
    ]
    }
    """
