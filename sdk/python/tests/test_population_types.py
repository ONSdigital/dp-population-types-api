from dp_population_types_api_sdk_python.client import PopulationTypesAPIClient
from dp_population_types_api_sdk_python.models import (
    GetAreaResponse,
    GetAreasResponse,
    GetBlockedAreaCountResult,
    GetCategorisationsResponse,
    GetDimensionCategoriesResponse,
    GetDimensionsResponse,
    GetPopulationTypeMetadataResponse,
    GetPopulationTypeResponse,
    GetPopulationTypesResponse,
    HTTPHeaders,
)

from .fakes import FakeResponse, FakeSession

BASE_URL = "http://localhost:27300"


def make_client(payload: dict) -> tuple[PopulationTypesAPIClient, FakeSession]:
    session = FakeSession(FakeResponse(200, payload))
    client = PopulationTypesAPIClient(base_url=BASE_URL, session=session)
    return client, session


# Headers


def test_headers_are_sent_when_provided() -> None:
    client, session = make_client({"population_type": {"name": "UR"}})

    client.population_types.get_population_type(
        "UR",
        headers=HTTPHeaders(access_token="token", florence_token="florence"),
    )

    assert session.last_kwargs["headers"] == {
        "Authorization": "Bearer token",
        "X-Florence-Token": "florence",
    }


def test_headers_are_none_when_not_provided() -> None:
    client, session = make_client({"population_type": {"name": "UR"}})

    client.population_types.get_population_type("UR")

    assert session.last_kwargs["headers"] is None


# Population types


def test_get_population_types_parses_response() -> None:
    payload = {
        "items": [
            {
                "name": "UR",
                "label": "Usual residents",
                "description": "All usual residents",
                "type": "microdata",
            }
        ],
        "count": 1,
        "offset": 0,
        "limit": 20,
        "total_count": 1,
    }
    client, session = make_client(payload)

    result = client.population_types.get_population_types(limit=20, offset=0)

    assert isinstance(result, GetPopulationTypesResponse)
    assert result.items[0].name == "UR"
    assert result.items[0].type == "microdata"
    assert result.total_count == 1
    assert session.last_kwargs["method"] == "GET"
    assert session.last_kwargs["url"] == f"{BASE_URL}/population-types"
    assert session.last_kwargs["params"] == {"offset": 0, "limit": 20}


def test_get_population_types_sends_require_default_dataset() -> None:
    client, session = make_client({"items": []})

    client.population_types.get_population_types(require_default_dataset=True)

    assert session.last_kwargs["params"] == {"require-default-dataset": "true"}


def test_get_population_type_parses_response() -> None:
    payload = {
        "population_type": {
            "name": "UR",
            "label": "Usual residents",
            "description": "All usual residents",
            "type": "microdata",
        }
    }
    client, session = make_client(payload)

    result = client.population_types.get_population_type("UR")

    assert isinstance(result, GetPopulationTypeResponse)
    assert result.population_type is not None
    assert result.population_type.label == "Usual residents"
    assert session.last_kwargs["method"] == "GET"
    assert session.last_kwargs["url"] == f"{BASE_URL}/population-types/UR"


def test_get_population_type_metadata_parses_response() -> None:
    payload = {
        "population_type": "UR",
        "default_dataset_id": "TS008",
        "edition": "2021",
        "version": 1,
    }
    client, session = make_client(payload)

    result = client.population_types.get_population_type_metadata("UR")

    assert isinstance(result, GetPopulationTypeMetadataResponse)
    assert result.default_dataset_id == "TS008"
    assert result.version == 1
    assert session.last_kwargs["method"] == "GET"
    assert session.last_kwargs["url"] == f"{BASE_URL}/population-types/UR/metadata"


# Areas


def test_get_area_parses_response() -> None:
    payload = {"area": {"id": "E06000001", "label": "Hartlepool", "area_type": "ltla"}}
    client, session = make_client(payload)

    result = client.population_types.get_area("UR", "ltla", "E06000001")

    assert isinstance(result, GetAreaResponse)
    assert result.area is not None
    assert result.area.label == "Hartlepool"
    assert result.area.area_type == "ltla"
    assert session.last_kwargs["method"] == "GET"
    assert session.last_kwargs["url"] == (
        f"{BASE_URL}/population-types/UR/area-types/ltla/areas/E06000001"
    )


def test_get_areas_parses_response() -> None:
    payload = {
        "items": [{"id": "E06000001", "label": "Hartlepool", "area_type": "ltla"}],
        "count": 1,
        "offset": 0,
        "limit": 10,
        "total_count": 1,
    }
    client, session = make_client(payload)

    result = client.population_types.get_areas(
        "UR", "ltla", text="Hart", limit=10, offset=0
    )

    assert isinstance(result, GetAreasResponse)
    assert result.items[0].id == "E06000001"
    assert result.count == 1
    assert session.last_kwargs["method"] == "GET"
    assert session.last_kwargs["url"] == (
        f"{BASE_URL}/population-types/UR/area-types/ltla/areas"
    )
    assert session.last_kwargs["params"] == {"offset": 0, "limit": 10, "q": "Hart"}


def test_get_areas_omits_unset_params() -> None:
    client, session = make_client({"items": []})

    client.population_types.get_areas("UR", "ltla")

    assert session.last_kwargs["params"] == {}


def test_get_blocked_area_count_parses_response() -> None:
    payload = {"passed": 10, "blocked": 2, "total": 12, "table_error": "too many"}
    client, session = make_client(payload)

    result = client.population_types.get_blocked_area_count(
        "UR",
        variables=["ltla", "sex"],
        filter_variable="ltla",
        filter_codes=["E06000001", "E06000002"],
    )

    assert isinstance(result, GetBlockedAreaCountResult)
    assert result.passed == 10
    assert result.blocked == 2
    assert result.total == 12
    assert result.table_error == "too many"
    assert session.last_kwargs["method"] == "GET"
    assert session.last_kwargs["url"] == (
        f"{BASE_URL}/population-types/UR/blocked-areas-count"
    )
    assert session.last_kwargs["params"] == {
        "vars": "ltla,sex",
        "fvar": "ltla",
        "areas": "E06000001,E06000002",
    }


def test_get_blocked_area_count_without_filter() -> None:
    client, session = make_client({"passed": 1, "blocked": 0, "total": 1})

    result = client.population_types.get_blocked_area_count("UR", variables=["ltla"])

    assert result.table_error is None
    assert session.last_kwargs["params"] == {"vars": "ltla"}


# Dimensions


def test_get_dimension_categories_parses_response() -> None:
    payload = {
        "items": [
            {
                "id": "sex",
                "label": "Sex",
                "quality_statement_text": "quality",
                "categories": [
                    {"id": "1", "label": "Female"},
                    {"id": "2", "label": "Male"},
                ],
            }
        ],
        "count": 1,
        "offset": 0,
        "limit": 20,
        "total_count": 1,
    }
    client, session = make_client(payload)

    result = client.population_types.get_dimension_categories(
        "UR", dimensions=["sex", "resident_age_7b"], limit=20, offset=0
    )

    assert isinstance(result, GetDimensionCategoriesResponse)
    assert result.items[0].id == "sex"
    assert result.items[0].categories[1].label == "Male"
    assert session.last_kwargs["method"] == "GET"
    assert session.last_kwargs["url"] == (
        f"{BASE_URL}/population-types/UR/dimension-categories"
    )
    assert session.last_kwargs["params"] == {
        "offset": 0,
        "limit": 20,
        "dims": "sex,resident_age_7b",
    }


def test_get_dimensions_description_parses_response() -> None:
    payload = {
        "items": [
            {
                "id": "sex",
                "label": "Sex",
                "description": "The classification of a person",
                "total_count": 2,
                "quality_statement_text": "quality",
            }
        ],
        "count": 1,
        "total_count": 1,
    }
    client, session = make_client(payload)

    result = client.population_types.get_dimensions_description(
        "UR", dimension_ids=["sex", "resident_age_7b"]
    )

    assert isinstance(result, GetDimensionsResponse)
    assert result.items[0].description == "The classification of a person"
    assert result.items[0].total_count == 2
    assert session.last_kwargs["method"] == "GET"
    assert session.last_kwargs["url"] == (
        f"{BASE_URL}/population-types/UR/dimensions-description"
    )
    assert session.last_kwargs["params"] == {"q": ["sex", "resident_age_7b"]}


def test_get_dimensions_description_without_ids() -> None:
    client, session = make_client({"items": []})

    client.population_types.get_dimensions_description("UR")

    assert session.last_kwargs["params"] is None


def test_get_categorisations_parses_response() -> None:
    payload = {
        "items": [
            {
                "id": "resident_age_7b",
                "label": "Age (7 categories)",
                "quality_statement_text": "quality",
                "default_categorisation": True,
                "categories": [{"id": "1", "label": "Aged 15 years and under"}],
            }
        ],
        "count": 1,
        "offset": 0,
        "limit": 20,
        "total_count": 5,
    }
    client, session = make_client(payload)

    result = client.population_types.get_categorisations(
        "UR", "resident_age_7b", limit=20, offset=0
    )

    assert isinstance(result, GetCategorisationsResponse)
    assert result.items[0].default_categorisation is True
    assert result.items[0].categories[0].label == "Aged 15 years and under"
    assert result.total_count == 5
    assert session.last_kwargs["method"] == "GET"
    assert session.last_kwargs["url"] == (
        f"{BASE_URL}/population-types/UR/dimensions/resident_age_7b/categorisations"
    )
    assert session.last_kwargs["params"] == {"offset": 0, "limit": 20}
