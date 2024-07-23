import datetime
from typing import List

from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class OrganizationRequest(GenericRequest):
    pass


class OrganizationListRequest(GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class OrganizationResponse(GenericResponse):
    name: str
    create_time: datetime.datetime
    update_time: datetime.datetime


class OrganizationListResponse(GenericListResponse):
    organizations: List[OrganizationResponse]
