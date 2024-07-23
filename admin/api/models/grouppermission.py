import datetime
from typing import List

from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class GroupPermissionRequest(GenericRequest):
    pass


class GroupPermissionListRequest(GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class GroupPermissionResponse(GenericResponse):
    group_id: str
    user_id: str
    create_time: datetime.datetime


class GroupPermissionListResponse(GenericListResponse):
    grouppermissions: List[GroupPermissionResponse]
