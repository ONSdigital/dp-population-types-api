Feature: Filter Metadata
  As an API user

  I want to save some metadata against a dataset

  Background:

  Given private endpoints are enabled

  And I am identified as "user@ons.gov.uk"

  And I am authorised

  Scenario: The root population-types endpoint should save metadata
    When I PUT "/population-types/UR/metadata"
    """
    {
      "default-dataset-id": "default-dataset"
    }
    """
    Then the HTTP status code should be "201"
    And a document in collection "filterMetadata" with key "id" value "UR" should match:
    """
    {
      "id": "UR",
      "default-dataset-id": "default-dataset"
    }
    """
  Scenario: The root population-types endpoint should get metadata
    Given I have this metadata:
    """
    [
      {
        "_id": "UR",
        "default-dataset-id": "default-dataset",
        "id": "UR"
      }
    ]
    """
    When I GET "/population-types/UR/metadata"
    Then the HTTP status code should be "201"
    And I should receive the following JSON response:
    """
    {
      "population-type": "UR",
      "default-dataset-id": "default-dataset"
    }
    """

  Scenario: The root population-types endpoint should return 404 if Blob not found
    When I PUT "/population-types/DOESNT-EXIST/metadata"
    """
    {
      "default-dataset-id": "default-dataset"
    }
    """
    Then the HTTP status code should be "404"
