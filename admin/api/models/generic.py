import datetime
from typing import List

from fastapi import status, HTTPException
from pydantic import BaseModel, Field


# --- REQUEST MODELS --- #
class GenericRequest(BaseModel):
    id: str = Field(...)


class GenericPaginationRequest(BaseModel):
    page: int | None = Field(default=1)
    size: int | None = Field(default=10)


# --- RESPONSE MODELS --- #
class GenericResponse(BaseModel):
    id: str


class GenericListResponse(BaseModel):
    totalElements: int
    page: int
    size: int
    data: List[None]


class GenericNotFoundResponse(BaseModel):
    status_code: int = status.HTTP_404_NOT_FOUND
    detail: str = 'Not found'


class GenericUnauthorizedResponse(BaseModel):
    status_code: int = status.HTTP_401_UNAUTHORIZED
    detail: str = 'Unauthorized'


# --- TOKEN SPECIFIC --- #
class GenericTokenPayload(BaseModel):
    iat: datetime.datetime
    iss: str
    aud: str
    exp: datetime.datetime


class GenericTokenRequest(BaseModel):
    pass


class GenericTokenResponse(BaseModel):
    pass
