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

from fastapi import APIRouter, Depends, status

from ..database import fetch_user_organizations, fetch_user_workspaces, fetch_user_groups
from ..errors import NotFoundError, EmptyDataException, NoContentError, NotFoundException, \
    UnknownApiError
from ..models import UserOrganizationListRequest, UserOrganizationListResponse, \
    UserWorkspaceListResponse, UserWorkspaceListRequest, UserGroupListResponse, UserGroupListRequest

users_api_router = APIRouter(
    prefix='/user',
    tags=['user'],
)


# --- GET --- #
@users_api_router.get(path="/organizations",
                      responses={
                          status.HTTP_200_OK: {
                              'model': UserOrganizationListResponse
                          }
                      }
                      )
async def get_user_organizations(data: Annotated[UserOrganizationListRequest, Depends()]):
    try:
        organizations, count = fetch_user_organizations(user_id=data.id, page=data.page, size=data.size)

        return UserOrganizationListResponse(data=organizations,
                                            totalElements=count,
                                            page=data.page,
                                            size=data.size)
    except EmptyDataException:
        return NoContentError()
    except NotFoundException as e:
        return NotFoundError(message=str(e))
    except Exception as e:
        return UnknownApiError(message=str(e))


@users_api_router.get(path="/workspaces",
                      responses={
                          status.HTTP_200_OK: {
                              'model': UserWorkspaceListResponse
                          }
                      }
                      )
async def get_user_workspaces(data: Annotated[UserWorkspaceListRequest, Depends()]):
    try:
        workspaces, count = fetch_user_workspaces(user_id=data.id, page=data.page, size=data.size)

        return UserWorkspaceListResponse(data=workspaces,
                                         totalElements=count,
                                         page=data.page,
                                         size=data.size)
    except EmptyDataException:
        return NoContentError()
    except NotFoundException as e:
        return NotFoundError(message=str(e))
    except Exception as e:
        return UnknownApiError(message=str(e))


@users_api_router.get(path="/groups",
                      responses={
                          status.HTTP_200_OK: {
                              'model': UserGroupListResponse
                          }
                      }
                      )
async def get_user_groups(data: Annotated[UserGroupListRequest, Depends()]):
    try:
        groups, count = fetch_user_groups(user_id=data.id, page=data.page, size=data.size)

        return UserGroupListResponse(data=groups,
                                     totalElements=count,
                                     page=data.page,
                                     size=data.size)
    except EmptyDataException:
        return NoContentError()
    except NotFoundException as e:
        return NotFoundError(message=str(e))
    except Exception as e:
        return UnknownApiError(message=str(e))

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
