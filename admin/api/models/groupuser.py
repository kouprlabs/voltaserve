import datetime
from typing import List

from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class GroupUserRequest(GenericRequest):
    pass


class GroupUserListRequest(GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class GroupUserResponse(GenericResponse):
    group_id: str
    user_id: str
    create_time: datetime.datetime


class GroupUserListResponse(GenericListResponse):
    data: List[GroupUserResponse]
