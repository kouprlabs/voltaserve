from typing import List

from .generic import GenericPaginationRequest, GenericListResponse, GenericRequest, BaseModel


# --- REQUEST MODELS --- #
class IndexRequest(GenericRequest):
    pass


class IndexListRequest(GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class IndexResponse(BaseModel):
    tablename: str
    indexname: str
    indexdef: str


class IndexListResponse(GenericListResponse):
    data: List[IndexResponse]
