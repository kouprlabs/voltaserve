from pydantic import BaseModel, Field


# --- REQUEST MODELS --- #
class GenericRequest(BaseModel):
    id: str = Field(...)


class GenericPaginationRequest(BaseModel):
    page: int | None = Field(default=0)
    size: int | None = Field(default=10)


# --- RESPONSE MODELS --- #
class GenericResponse(BaseModel):
    id: str


class GenericListResponse(BaseModel):
    pass
