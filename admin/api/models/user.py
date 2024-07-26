import datetime
from typing import List

from pydantic import EmailStr

from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class UserRequest(GenericRequest):
    pass


class UserListRequest(GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class UserResponse(GenericResponse):
    full_name: str
    username: str
    email: EmailStr
    is_email_confirmed: bool
    picture: str | None
    create_time: datetime.datetime
    update_time: datetime.datetime


class UserListResponse(GenericListResponse):
    data: List[UserResponse]
