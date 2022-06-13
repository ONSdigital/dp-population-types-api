Feature: Area Types

  Scenario: Getting area-types happy
	When a geography query response is available from Cantabular api extension
	
	And I GET "/population-types/Example/area-types"
	
	Then I should receive the following JSON response:
	"""
	{
		"area-types":[
		  {
	  		"id":"country",
				"label":"Country",
				"total_count": 2
		  },
		  {
				"id":"city",
				"label":"City",
				"total_count": 3
		  }
		]
	}
	"""
	
	And the HTTP status code should be "200"

  Scenario: Getting area-types not found
	When an error json response is returned from Cantabular api extension
	
	And I GET "/population-types/Inexistent/area-types"
	
	Then I should receive the following JSON response:
	"""
	{"errors":["failed to get area-types: error(s) returned by graphQL query"]}
	"""
	
	And the HTTP status code should be "404"
