from pydantic import BaseModel, Field

from .pagination import PaginationResponse

# Population types


class PopulationType(BaseModel):
    name: str | None = None
    label: str | None = None
    description: str | None = None
    type: str | None = None


class GetPopulationTypeResponse(BaseModel):
    population_type: PopulationType | None = None


class GetPopulationTypesResponse(PaginationResponse):
    items: list[PopulationType] = Field(default_factory=list)


class GetPopulationTypeMetadataResponse(BaseModel):
    population_type: str | None = None
    default_dataset_id: str | None = None
    edition: str | None = None
    version: int | None = None


# Areas


class Area(BaseModel):
    id: str | None = None
    label: str | None = None
    area_type: str | None = None


class GetAreaResponse(BaseModel):
    area: Area | None = None


class GetAreasResponse(PaginationResponse):
    items: list[Area] = Field(default_factory=list)


class GetBlockedAreaCountResult(BaseModel):
    passed: int = 0
    blocked: int = 0
    total: int = 0
    table_error: str | None = None


# Dimensions


class Dimension(BaseModel):
    id: str | None = None
    label: str | None = None
    description: str | None = None
    total_count: int = 0
    quality_statement_text: str | None = None


class GetDimensionsResponse(PaginationResponse):
    items: list[Dimension] = Field(default_factory=list)


class DimensionCategory(BaseModel):
    """A single category (option) within a dimension"""

    id: str | None = None
    label: str | None = None


class Category(BaseModel):
    """A dimension/categorisation along with its categories"""

    id: str | None = None
    label: str | None = None
    categories: list[DimensionCategory] = Field(default_factory=list)
    quality_statement_text: str | None = None
    default_categorisation: bool | None = None


class GetDimensionCategoriesResponse(PaginationResponse):
    items: list[Category] = Field(default_factory=list)


class GetCategorisationsResponse(PaginationResponse):
    items: list[Category] = Field(default_factory=list)
