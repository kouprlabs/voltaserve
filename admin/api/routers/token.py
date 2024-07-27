# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

from typing import Annotated

from fastapi import APIRouter, status, Header, Depends

from ..dependencies import JWTBearer
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
    },
    dependencies=[Depends(JWTBearer())]
)


# --- GET --- #

# --- PATCH --- #

# --- POST --- #
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

# --- PUT --- #

# --- DELETE --- #
