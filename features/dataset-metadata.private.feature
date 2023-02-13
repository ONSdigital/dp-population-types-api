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
      "default_dataset_id": "default-dataset",
      "edition": "2021",
      "version": 2
    }
    """
    Then the HTTP status code should be "201"
    And a document in collection "defaultDatasetMetadata" with key "id" value "UR" should match:
    """
    {
      "_id": "UR",
      "id": "UR",
      "default_dataset_id": "default-dataset",
      "edition": "2021",
      "version": {
          "$numberInt":"2"
        }
    }
    """

  Scenario: The root population-types endpoint should get metadata
    Given I have this metadata:
    """
    [
      {
        "_id": "UR",
        "default_dataset_id": "default-dataset",
        "id": "UR",
        "edition": "2021",
        "version": 2
      }
    ]
    """
    When I GET "/population-types/UR/metadata"
    Then the HTTP status code should be "200"
    And I should receive the following JSON response:
    """
    {
      "population_type": "UR",
      "default_dataset_id": "default-dataset",
      "edition": "2021",
      "version": 2
    }
    """

  Scenario: The root population-types endpoint should return 404 if Blob not found
    When I GET "/population-types/DOESNT-EXIST/metadata"
    """
    {
      "default_dataset_id": "default-dataset"
    }
    """
    Then the HTTP status code should be "404"
