# dp-population-types-api-sdk-python

Python SDK for interacting with `dp-population-types-api`.

## Overview

This SDK provides a Python client for interacting with `dp-population-types-api`. It is intended to be consumed by services that require endpoints from the dp-population-types-api, such as `frontend-dataset-controller`. Responses are returned as validated Pydantic models.

## Available client methods

| Name | Description |
| ------ | ------------- |
| [`health`](#health) | Returns the health status of the API |
| [`get_population_types`](#get_population_types) | Returns a paginated list of population types |
| [`get_population_type`](#get_population_type) | Returns a single population type |
| [`get_population_type_metadata`](#get_population_type_metadata) | Returns the default dataset metadata for a population type |
| [`get_area`](#get_area) | Returns a single area for a given population type and area type |
| [`get_areas`](#get_areas) | Returns a paginated list of areas for a given population type and area type, optionally filtered by search text |
| [`get_blocked_area_count`](#get_blocked_area_count) | Returns the passed, blocked and total area counts for a combination of variables |
| [`get_dimension_categories`](#get_dimension_categories) | Returns the categories for the given dimensions |
| [`get_dimensions_description`](#get_dimensions_description) | Returns the descriptions for the given dimensions |
| [`get_categorisations`](#get_categorisations) | Returns a paginated list of categorisations for a dimension |

All methods except `health` are on the `population_types` resource, e.g. `client.population_types.get_areas(...)`.

## Requirements

- Python `>=3.14`
- `requests`
- `pydantic`

## Install

### Install in this project with Makefile

```bash
make install-dev
make install
```

### Install from Git in another service

Replace `<release-tag>` with a tag from the [dp-population-types-api releases](https://github.com/ONSdigital/dp-population-types-api/releases) page.

```bash
pip install "git+https://github.com/ONSdigital/dp-population-types-api.git@<release-tag>#subdirectory=sdk/python"
```

## Instantiation

Example using `create_client`:

```python
from dp_population_types_api_sdk_python import create_client

client = create_client(base_url="http://localhost:27300")
```

Example using a shared `requests.Session`, so requests reuse connections and can share default headers:

```python
import requests

from dp_population_types_api_sdk_python import create_client

session = requests.Session()

client = create_client(
    base_url="http://localhost:27300",
    timeout=5.0,
    session=session,
)
```

## Example usage of client

This example demonstrates how `get_population_type()` could be used:

```python
from dp_population_types_api_sdk_python import (
    HTTPHeaders,
    NotFoundError,
    create_client,
)

client = create_client(base_url="http://localhost:27300")

headers = HTTPHeaders(access_token="example-auth-token")

try:
    response = client.population_types.get_population_type("UR", headers=headers)
except NotFoundError as exc:

    print(exc.status_code, exc)

```

## Available functionality

Every method on `population_types` accepts an optional `headers` argument. Optional parameters that are `None` are not sent, so the API defaults apply.

### health

```python
health = client.health()
```

Returns: `dict`

### get_population_types

```python
response = client.population_types.get_population_types(
    require_default_dataset=True,
    limit=20,
    offset=0,
)
```

Returns: `GetPopulationTypesResponse`

### get_population_type

```python
response = client.population_types.get_population_type("UR")
```

Returns: `GetPopulationTypeResponse`

### get_population_type_metadata

```python
response = client.population_types.get_population_type_metadata("UR")
```

Returns: `GetPopulationTypeMetadataResponse`

### get_area

```python
response = client.population_types.get_area("UR", "ltla", "E06000001")
```

Returns: `GetAreaResponse`

### get_areas

```python
response = client.population_types.get_areas(
    "UR",
    "ltla",
    text="Hart",
    limit=20,
    offset=0,
)
```

Returns: `GetAreasResponse`

### get_blocked_area_count

```python
response = client.population_types.get_blocked_area_count(
    "UR",
    variables=["ltla", "sex"],
    filter_variable="ltla",
    filter_codes=["E06000001", "E06000002"],
)
```

Returns: `GetBlockedAreaCountResult`

### get_dimension_categories

```python
response = client.population_types.get_dimension_categories(
    "UR",
    dimensions=["sex", "resident_age_7b"],
    limit=20,
    offset=0,
)
```

Returns: `GetDimensionCategoriesResponse`

### get_dimensions_description

```python
response = client.population_types.get_dimensions_description(
    "UR",
    dimension_ids=["sex", "resident_age_7b"],
)
```

Returns: `GetDimensionsResponse`

### get_categorisations

```python
response = client.population_types.get_categorisations(
    "UR",
    "resident_age_7b",
    limit=20,
    offset=0,
)
```

Returns: `GetCategorisationsResponse`

## Additional Information

### Errors

The client raises typed exceptions for unsuccessful responses. Each exception has a `status_code` attribute, and its message contains the response body returned by the API where there is one (e.g. `{"errors": ["population type not found"]}`).

| Exception | Raised when |
| --------- | ----------- |
| `AuthenticationError` | The API returns `401` or `403` |
| `NotFoundError` | The API returns `404` |
| `ValidationError` | The API returns `400` or `422` |
| `APIError` | Any other error response, or the request fails to send (e.g. connection error). Base class of all the above |

```python
from dp_population_types_api_sdk_python import (
    APIError,
    AuthenticationError,
    NotFoundError,
    ValidationError,
    create_client,
)

client = create_client(base_url="http://localhost:27300")

try:
    response = client.population_types.get_population_type("UR")
except NotFoundError:
    print("Population type does not exist")
except AuthenticationError:
    print("Authentication failed")
except ValidationError as exc:
    print(f"Validation failed: {exc}")
except APIError as exc:
    print(f"API failed: status={exc.status_code} message={exc}")
```

## Development commands

```bash
make test
make lint
make mypy   
make format
make audit
```