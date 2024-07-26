import datetime
from typing import List

from pydantic import EmailStr

from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class InvitationRequest(GenericRequest):
    pass


class InvitationListRequest(GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class InvitationResponse(GenericResponse):
    organization_id: str
    owner_id: str
    email: EmailStr
    status: str
    create_time: datetime.datetime
    update_time: datetime.datetime


class InvitationListResponse(GenericListResponse):
    data: List[InvitationResponse]
