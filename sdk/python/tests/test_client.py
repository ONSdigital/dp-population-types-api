import pytest
import requests

from dp_population_types_api_sdk_python import create_client
from dp_population_types_api_sdk_python.client import PopulationTypesAPIClient
from dp_population_types_api_sdk_python.exceptions import (
    APIError,
    AuthenticationError,
    NotFoundError,
    ValidationError,
)
from dp_population_types_api_sdk_python.protocols import (
    PopulationTypesAPIClientProtocol,
)

from .fakes import FakeResponse, FakeSession


def test_init_removes_trailing_slash_from_base_url() -> None:
    client = PopulationTypesAPIClient("http://localhost:27300/")

    assert client.base_url == "http://localhost:27300"


def test_init_uses_supplied_session() -> None:
    session = FakeSession(FakeResponse(200, {"status": "ok"}))

    client = PopulationTypesAPIClient(
        base_url="http://localhost:27300",
        session=session,
    )

    assert client.session is session


def test_init_rejects_non_session_objects() -> None:
    with pytest.raises(TypeError):
        PopulationTypesAPIClient(
            base_url="http://localhost:27300",
            session=object(),
        )


def test_request_forwards_arguments_to_session() -> None:
    response = FakeResponse(200, {"result": "ok"})
    session = FakeSession(response)
    client = PopulationTypesAPIClient(
        base_url="http://localhost:27300",
        timeout=10,
        session=session,
    )

    result = client._request(
        "POST",
        "/endpoint",
        params={"offset": 1},
        json={"field": "value"},
        headers={"example-header": "header-value"},
    )

    assert result is response
    assert session.last_kwargs == {
        "method": "POST",
        "url": "http://localhost:27300/endpoint",
        "params": {"offset": 1},
        "json": {"field": "value"},
        "timeout": 10,
        "headers": {"example-header": "header-value"},
    }


def test_health_returns_json_and_requests_health_endpoint() -> None:
    session = FakeSession(FakeResponse(200, {"status": "ok"}))
    client = PopulationTypesAPIClient(
        base_url="http://localhost:27300",
        session=session,
    )

    assert client.health() == {"status": "ok"}
    assert session.last_kwargs["method"] == "GET"
    assert session.last_kwargs["url"] == "http://localhost:27300/health"


def test_request_wraps_network_errors() -> None:
    session = FakeSession(error=requests.ConnectionError("connection failed"))
    client = PopulationTypesAPIClient(
        base_url="http://localhost:27300",
        session=session,
    )

    with pytest.raises(APIError):
        client._request("GET", "/health")


@pytest.mark.parametrize(
    ("status_code", "exception_type"),
    [
        (400, ValidationError),
        (401, AuthenticationError),
        (403, AuthenticationError),
        (404, NotFoundError),
        (422, ValidationError),
        (429, APIError),
        (500, APIError),
        (502, APIError),
    ],
)
def test_request_maps_error_status_codes(
    status_code: int,
    exception_type: type[APIError],
) -> None:
    session = FakeSession(FakeResponse(status_code, {"errors": ["request failed"]}))
    client = PopulationTypesAPIClient(
        base_url="http://localhost:27300",
        session=session,
    )

    with pytest.raises(exception_type) as exc_info:
        client._request("GET", "/resource")

    assert exc_info.value.status_code == status_code


def test_request_error_includes_response_body() -> None:
    session = FakeSession(FakeResponse(404, {"errors": ["failed to get area"]}))
    client = PopulationTypesAPIClient(
        base_url="http://localhost:27300",
        session=session,
    )

    with pytest.raises(NotFoundError) as exc_info:
        client._request("GET", "/resource")

    assert "failed to get area" in str(exc_info.value)


def test_create_client_returns_protocol_conforming_client() -> None:
    client = create_client("http://localhost:27300")

    assert isinstance(client, PopulationTypesAPIClientProtocol)
    assert isinstance(client, PopulationTypesAPIClient)
