from typing import Annotated

from fastapi import APIRouter, Depends, status, Header

from ..exceptions import GenericUnauthorizedException
from ..models import GenericNotFoundResponse, TokenResponse, GenericUnauthorizedResponse

token_api_router = APIRouter(
    prefix='/token',
    tags=['token'],
    responses={
        status.HTTP_404_NOT_FOUND: {
            'model': GenericNotFoundResponse
        },
        status.HTTP_401_UNAUTHORIZED: {
            'model': GenericUnauthorizedResponse
        }
    }
)


# --- GET --- #
@token_api_router.post(path="/validate",
                       responses={
                           status.HTTP_200_OK: {
                               'model': TokenResponse
                           },
                       }
                       )
async def validate_token(authorization: Annotated[str | None, Header()] = None):
    if authorization is None:
        raise GenericUnauthorizedException()

    return TokenResponse(authorized=True)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
