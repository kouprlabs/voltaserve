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
from ..exceptions import GenericNotFoundException
from ..models import GenericNotFoundResponse, OrganizationUserListRequest, OrganizationUserListResponse, \
    WorkspaceUserListResponse, WorkspaceUserListRequest, GroupUserListResponse, GroupUserListRequest

users_api_router = APIRouter(
    prefix='/user',
    tags=['user'],
    responses={
        status.HTTP_404_NOT_FOUND: {
            'model': GenericNotFoundResponse
        }
    }
)


# --- GET --- #
@users_api_router.get(path="/organizations",
                      responses={
                          status.HTTP_200_OK: {
                              'model': OrganizationUserListResponse
                          }
                      }
                      )
async def get_user_organizations(data: Annotated[OrganizationUserListRequest, Depends()]):
    organizations, count = fetch_user_organizations(user_id=data.id, page=data.page, size=data.size)
    if organizations is None:
        raise GenericNotFoundException(detail='This user has no organizations')

    return OrganizationUserListResponse(data=organizations,
                                        totalElements=count,
                                        page=data.page,
                                        size=data.size)


@users_api_router.get(path="/workspaces",
                      responses={
                          status.HTTP_200_OK: {
                              'model': WorkspaceUserListResponse
                          }
                      }
                      )
async def get_user_workspaces(data: Annotated[WorkspaceUserListRequest, Depends()]):
    workspaces, count = fetch_user_workspaces(user_id=data.id, page=data.page, size=data.size)
    if workspaces is None:
        raise GenericNotFoundException(detail='This user has no workspaces')

    return WorkspaceUserListResponse(data=workspaces,
                                     totalElements=count,
                                     page=data.page,
                                     size=data.size)


@users_api_router.get(path="/groups",
                      responses={
                          status.HTTP_200_OK: {
                              'model': GroupUserListResponse
                          }
                      }
                      )
async def get_user_groups(data: Annotated[GroupUserListRequest, Depends()]):
    groups, count = fetch_user_groups(user_id=data.id, page=data.page, size=data.size)
    if groups is None:
        raise GenericNotFoundException(detail='This user has no groups')

    return GroupUserListResponse(data=groups,
                                 totalElements=count,
                                 page=data.page,
                                 size=data.size)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
