# dp-population-types-api-sdk-python

Python SDK for interacting with `dp-population-types-api`.

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

## Quick Start

```python
from dp_population_types_api_sdk_python import create_client

client = create_client(base_url="https://api.example.com")
health = client.health()

print(health)
```

## Population types

All operations are on the `population_types` resource and return validated Pydantic models.

```python
from dp_population_types_api_sdk_python import create_client

client = create_client(base_url="https://api.example.com")
pt = client.population_types

pt.get_population_types(require_default_dataset=True, limit=20, offset=0)
pt.get_population_type("UR")
pt.get_population_type_metadata("UR")

pt.get_area("UR", "ltla", "E06000001")
pt.get_areas("UR", "ltla", text="Hart", limit=20, offset=0)
pt.get_blocked_area_count(
    "UR",
    variables=["ltla", "sex"],
    filter_variable="ltla",
    filter_codes=["E06000001", "E06000002"],
)

pt.get_dimension_categories("UR", dimensions=["sex", "resident_age_7b"])
pt.get_dimensions_description("UR", dimension_ids=["sex", "resident_age_7b"])
pt.get_categorisations("UR", "resident_age_7b", limit=20, offset=0)
```

Optional parameters that are `None` are not sent, so the API defaults apply.

## Headers and authentication

You can pass per-request headers using `HTTPHeaders`. Pass the raw token; the `Bearer ` prefix is added for you.

```python
from dp_population_types_api_sdk_python import HTTPHeaders, create_client

client = create_client(base_url="https://api.example.com")

population_type = client.population_types.get_population_type(
    "UR",
    headers=HTTPHeaders(access_token="YOUR_TOKEN"),
)
```

| Field | Header sent |
|---|---|
| `access_token` | `Authorization: Bearer <access_token>` |
| `florence_token` | `X-Florence-Token: <florence_token>` |

`None` values are omitted before sending the request.

This SDK is designed to be used with a shared `requests.Session`. Pass one into `create_client()` so your requests reuse connections and can share default headers.

```python
import requests

from dp_population_types_api_sdk_python import create_client

session = requests.Session()
session.headers.update({"Authorization": "Bearer YOUR_TOKEN"})

client = create_client(
    base_url="https://api.example.com",
    session=session,
)
```

## Error handling

The client raises typed exceptions for common HTTP failures. The exception message contains the response body where one is returned.

```python
from dp_population_types_api_sdk_python import (
    APIError,
    AuthenticationError,
    NotFoundError,
    ValidationError,
    create_client,
)

client = create_client(base_url="https://api.example.com")

try:
    population_type = client.population_types.get_population_type("UR")
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