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

from ..database import fetch_workspace, fetch_workspaces, fetch_organization_workspaces
from ..dependencies import JWTBearer
from ..exceptions import GenericNotFoundException
from ..models import GenericNotFoundResponse, WorkspaceResponse, WorkspaceRequest, WorkspaceListResponse, \
    WorkspaceListRequest, OrganizationWorkspaceListRequest

workspace_api_router = APIRouter(
    prefix='/workspace',
    tags=['workspace'],
    responses={
        status.HTTP_404_NOT_FOUND: {
            'model': GenericNotFoundResponse
        }
    },
    dependencies=[Depends(JWTBearer())]
)


# --- GET --- #
@workspace_api_router.get(path="/",
                          responses={
                              status.HTTP_200_OK: {
                                  'model': WorkspaceResponse
                              }
                          })
async def get_workspace(data: Annotated[WorkspaceRequest, Depends()]):
    workspace = fetch_workspace(_id=data.id)
    if workspace is None:
        raise GenericNotFoundException(detail=f'Workspace with id={data.id} does not exist')

    return WorkspaceResponse(**workspace)


@workspace_api_router.get(path="/all",
                          responses={
                              status.HTTP_200_OK: {
                                  'model': WorkspaceListResponse
                              }
                          }
                          )
async def get_all_workspaces(data: Annotated[WorkspaceListRequest, Depends()]):
    workspaces, count = fetch_workspaces(page=data.page, size=data.size)
    if workspaces is None:
        raise GenericNotFoundException(detail='This instance has no workspaces')

    return WorkspaceListResponse(data=workspaces, totalElements=count['count'], page=data.page, size=data.size)


@workspace_api_router.get(path="/organization",
                          responses={
                              status.HTTP_200_OK: {
                                  'model': WorkspaceListResponse
                              }
                          }
                          )
async def get_organization_workspaces(data: Annotated[OrganizationWorkspaceListRequest, Depends()]):
    workspaces, count = fetch_organization_workspaces(organization_id=data.id, page=data.page, size=data.size)
    if workspaces is None:
        raise GenericNotFoundException(detail='This instance has no workspaces')

    return WorkspaceListResponse(data=workspaces, totalElements=count['count'], page=data.page, size=data.size)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
