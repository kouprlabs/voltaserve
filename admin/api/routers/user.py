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

from aiohttp import ClientSession
from fastapi import APIRouter, Depends, status, Header

from ..database import fetch_user, fetch_users, fetch_user_organizations
from ..dependencies import JWTBearer, settings
from ..exceptions import GenericNotFoundException
from ..models import GenericNotFoundResponse, UserOrganizationListRequest, UserOrganizationListResponse, \
    UserListRequest, UserListResponse, UserRequest, UserResponse

users_api_router = APIRouter(
    prefix='/user',
    tags=['user'],
    responses={
        status.HTTP_404_NOT_FOUND: {
            'model': GenericNotFoundResponse
        }
    },
    dependencies=[Depends(JWTBearer())]
)


# --- GET --- #
@users_api_router.get(path="/",
                      responses={
                          status.HTTP_200_OK: {
                              'model': UserResponse
                          }
                      }
                      )
async def get_user(data: Annotated[UserRequest, Depends()]):
    user = fetch_user(_id=data.id)
    if user is None:
        raise GenericNotFoundException(detail=f'User with id={data.id} does not exist')

    return UserResponse(**user)


@users_api_router.get(path="/all",
                      responses={
                          status.HTTP_200_OK: {
                              'model': UserListResponse
                          }
                      }
                      )
async def get_all_users(data: Annotated[UserListRequest, Depends()], x_authorization: Annotated[str, Header()]):
    async with ClientSession() as sess:
        async with sess.get(f"{settings.idp_url}/v2/user/all?page=1&size=5", headers={'Authorization': x_authorization}) as resp:
            if resp.status != 200:
                raise GenericNotFoundException(detail='Wrong user authorization token')

            users = await resp.json()

        async with sess.get(f"{settings.idp_url}/v2/user/count", headers={'Authorization': x_authorization}) as resp:
            if resp.status != 200:
                raise GenericNotFoundException(detail='Wrong user authorization token')

            count = await resp.json()

    if users is None:
        raise GenericNotFoundException(detail='This instance has no users')

    return UserListResponse(data=users, totalElements=count['count'], page=data.page, size=data.size)


@users_api_router.get(path="/organizations",
                      responses={
                          status.HTTP_200_OK: {
                              'model': UserOrganizationListResponse
                          }
                      }
                      )
async def get_user_organizations(data: Annotated[UserOrganizationListRequest, Depends()]):
    organizations, count = fetch_user_organizations(user_id=data.id, page=data.page, size=data.size)
    if organizations is None:
        raise GenericNotFoundException(detail='This user has no organizations')

    return UserOrganizationListResponse(data=organizations,
                                        totalElements=count['count'],
                                        page=data.page,
                                        size=data.size)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
