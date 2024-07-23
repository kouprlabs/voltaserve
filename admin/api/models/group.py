import datetime
from typing import List

from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class GroupRequest(GenericRequest):
    pass


class GroupListRequest(GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class GroupResponse(GenericResponse):
    name: str
    organization_id: str
    create_time: datetime.datetime
    update_time: datetime.datetime


class GroupListResponse(GenericListResponse):
    groups: List[GroupResponse]
