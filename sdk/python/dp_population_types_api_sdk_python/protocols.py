from collections.abc import Mapping
from typing import Any, Protocol, runtime_checkable

import requests

from .models import (
    GetAreaResponse,
    GetAreasResponse,
    GetBlockedAreaCountResult,
    GetCategorisationsResponse,
    GetDimensionCategoriesResponse,
    GetDimensionsResponse,
    GetPopulationTypeMetadataResponse,
    GetPopulationTypeResponse,
    GetPopulationTypesResponse,
)


@runtime_checkable
class Headers(Protocol):
    def to_http_headers(self) -> Mapping[str, str]: ...


class RequestingClient(Protocol):
    def _request(
        self,
        method: str,
        path: str,
        *,
        params: dict[str, Any] | None = None,
        json: dict[str, Any] | None = None,
        headers: Mapping[str, str] | None = None,
    ) -> requests.Response: ...


@runtime_checkable
class PopulationTypesClientProtocol(Protocol):
    """Protocol for population-type-specific operations."""

    def get_population_types(
        self,
        require_default_dataset: bool = False,
        limit: int | None = None,
        offset: int | None = None,
        headers: Headers | None = None,
    ) -> GetPopulationTypesResponse: ...

    def get_population_type(
        self,
        population_type: str,
        headers: Headers | None = None,
    ) -> GetPopulationTypeResponse: ...

    def get_population_type_metadata(
        self,
        population_type: str,
        headers: Headers | None = None,
    ) -> GetPopulationTypeMetadataResponse: ...

    def get_area(
        self,
        population_type: str,
        area_type: str,
        area_id: str,
        headers: Headers | None = None,
    ) -> GetAreaResponse: ...

    def get_areas(
        self,
        population_type: str,
        area_type: str,
        text: str | None = None,
        limit: int | None = None,
        offset: int | None = None,
        headers: Headers | None = None,
    ) -> GetAreasResponse: ...

    def get_blocked_area_count(
        self,
        population_type: str,
        variables: list[str],
        filter_variable: str | None = None,
        filter_codes: list[str] | None = None,
        headers: Headers | None = None,
    ) -> GetBlockedAreaCountResult: ...

    def get_dimension_categories(
        self,
        population_type: str,
        dimensions: list[str],
        limit: int | None = None,
        offset: int | None = None,
        headers: Headers | None = None,
    ) -> GetDimensionCategoriesResponse: ...

    def get_dimensions_description(
        self,
        population_type: str,
        dimension_ids: list[str] | None = None,
        headers: Headers | None = None,
    ) -> GetDimensionsResponse: ...

    def get_categorisations(
        self,
        population_type: str,
        dimension: str,
        limit: int | None = None,
        offset: int | None = None,
        headers: Headers | None = None,
    ) -> GetCategorisationsResponse: ...


@runtime_checkable
class PopulationTypesAPIClientProtocol(Protocol):
    """Protocol for the PopulationTypesAPIClient."""

    def health(self) -> dict[str, Any]: ...

    population_types: PopulationTypesClientProtocol
