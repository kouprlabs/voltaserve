import datetime
from typing import List

from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class OrganizationUserRequest(GenericRequest):
    pass


class OrganizationUserListRequest(GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class OrganizationUserResponse(GenericResponse):
    organization_id: str
    user_id: str
    create_time: datetime.datetime


class OrganizationUserListResponse(GenericListResponse):
    OrganizationUsers: List[OrganizationUserResponse]
