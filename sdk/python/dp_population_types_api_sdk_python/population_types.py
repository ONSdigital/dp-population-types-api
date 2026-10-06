from typing import Any

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
    PaginationParams,
)
from .protocols import Headers, RequestingClient


class PopulationTypesAPI:
    def __init__(self, client: RequestingClient) -> None:
        self._client = client

    def get_population_types(
        self,
        require_default_dataset: bool = False,
        limit: int | None = None,
        offset: int | None = None,
        headers: Headers | None = None,
    ) -> GetPopulationTypesResponse:
        """GET /population-types"""

        params: dict[str, Any] = PaginationParams(
            offset=offset, limit=limit
        ).to_query_params()
        if require_default_dataset:
            params["require-default-dataset"] = "true"

        response = self._client._request(
            "GET",
            "/population-types",
            params=params,
            headers=headers.to_http_headers() if headers else None,
        )

        return GetPopulationTypesResponse.model_validate(response.json())

    def get_population_type(
        self,
        population_type: str,
        headers: Headers | None = None,
    ) -> GetPopulationTypeResponse:
        """GET /population-types/{population_type}"""

        response = self._client._request(
            "GET",
            f"/population-types/{population_type}",
            headers=headers.to_http_headers() if headers else None,
        )

        return GetPopulationTypeResponse.model_validate(response.json())

    def get_population_type_metadata(
        self,
        population_type: str,
        headers: Headers | None = None,
    ) -> GetPopulationTypeMetadataResponse:
        """GET /population-types/{population_type}/metadata"""

        response = self._client._request(
            "GET",
            f"/population-types/{population_type}/metadata",
            headers=headers.to_http_headers() if headers else None,
        )

        return GetPopulationTypeMetadataResponse.model_validate(response.json())

    def get_area(
        self,
        population_type: str,
        area_type: str,
        area_id: str,
        headers: Headers | None = None,
    ) -> GetAreaResponse:
        """GET /population-types/{population_type}/area-types/{area_type}/areas/{area_id}"""

        response = self._client._request(
            "GET",
            f"/population-types/{population_type}/area-types/{area_type}/areas/{area_id}",
            headers=headers.to_http_headers() if headers else None,
        )

        return GetAreaResponse.model_validate(response.json())

    def get_areas(
        self,
        population_type: str,
        area_type: str,
        text: str | None = None,
        limit: int | None = None,
        offset: int | None = None,
        headers: Headers | None = None,
    ) -> GetAreasResponse:
        """GET /population-types/{population_type}/area-types/{area_type}/areas

        text is sent as the "q" query param to search areas by name
        """

        params: dict[str, Any] = PaginationParams(
            offset=offset, limit=limit
        ).to_query_params()
        if text is not None:
            params["q"] = text

        response = self._client._request(
            "GET",
            f"/population-types/{population_type}/area-types/{area_type}/areas",
            params=params,
            headers=headers.to_http_headers() if headers else None,
        )

        return GetAreasResponse.model_validate(response.json())

    def get_blocked_area_count(
        self,
        population_type: str,
        variables: list[str],
        filter_variable: str | None = None,
        filter_codes: list[str] | None = None,
        headers: Headers | None = None,
    ) -> GetBlockedAreaCountResult:
        """GET /population-types/{population_type}/blocked-areas-count

        The API only applies the area filter when filter_variable is set
        """

        params: dict[str, Any] = {"vars": ",".join(variables)}
        if filter_variable is not None:
            params["fvar"] = filter_variable
        if filter_codes is not None:
            params["areas"] = ",".join(filter_codes)

        response = self._client._request(
            "GET",
            f"/population-types/{population_type}/blocked-areas-count",
            params=params,
            headers=headers.to_http_headers() if headers else None,
        )

        return GetBlockedAreaCountResult.model_validate(response.json())

    def get_dimension_categories(
        self,
        population_type: str,
        dimensions: list[str],
        limit: int | None = None,
        offset: int | None = None,
        headers: Headers | None = None,
    ) -> GetDimensionCategoriesResponse:
        """GET /population-types/{population_type}/dimension-categories"""

        params: dict[str, Any] = PaginationParams(
            offset=offset, limit=limit
        ).to_query_params()
        params["dims"] = ",".join(dimensions)

        response = self._client._request(
            "GET",
            f"/population-types/{population_type}/dimension-categories",
            params=params,
            headers=headers.to_http_headers() if headers else None,
        )

        return GetDimensionCategoriesResponse.model_validate(response.json())

    def get_dimensions_description(
        self,
        population_type: str,
        dimension_ids: list[str] | None = None,
        headers: Headers | None = None,
    ) -> GetDimensionsResponse:
        """GET /population-types/{population_type}/dimensions-description

        Dimension IDs are sent as repeated query params (?q=a&q=b), matching
        what the API handler expects
        """

        response = self._client._request(
            "GET",
            f"/population-types/{population_type}/dimensions-description",
            params={"q": dimension_ids} if dimension_ids is not None else None,
            headers=headers.to_http_headers() if headers else None,
        )

        return GetDimensionsResponse.model_validate(response.json())

    def get_categorisations(
        self,
        population_type: str,
        dimension: str,
        limit: int | None = None,
        offset: int | None = None,
        headers: Headers | None = None,
    ) -> GetCategorisationsResponse:
        """GET /population-types/{population_type}/dimensions/{dimension}/categorisations"""

        response = self._client._request(
            "GET",
            f"/population-types/{population_type}/dimensions/{dimension}/categorisations",
            params=PaginationParams(offset=offset, limit=limit).to_query_params(),
            headers=headers.to_http_headers() if headers else None,
        )

        return GetCategorisationsResponse.model_validate(response.json())
