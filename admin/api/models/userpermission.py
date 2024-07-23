import datetime
from typing import List

from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class UserPermissionRequest(GenericRequest):
    pass


class UserPermissionListRequest(GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class UserPermissionResponse(GenericResponse):
    user_id: str
    resource_id: str
    permission: str
    create_time: datetime.datetime


class UserPermissionListResponse(GenericListResponse):
    userpermissions: List[UserPermissionResponse]
