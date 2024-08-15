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

from fastapi import APIRouter, Depends, status, Response

from ..database import fetch_workspace, fetch_workspaces, update_workspace
from ..exceptions import GenericNotFoundException, GenericApiException
from ..models import GenericNotFoundResponse, WorkspaceResponse, WorkspaceRequest, WorkspaceListResponse, \
    WorkspaceListRequest, UpdateWorkspaceRequest, GenericUnexpectedErrorResponse, GenericAcceptedResponse

workspace_api_router = APIRouter(
    prefix='/workspace',
    tags=['workspace'],
    responses={
        status.HTTP_404_NOT_FOUND: {
            'model': GenericNotFoundResponse
        },
        status.HTTP_500_INTERNAL_SERVER_ERROR: {
            'model': GenericUnexpectedErrorResponse
        },
        status.HTTP_202_ACCEPTED: {
            'model': GenericAcceptedResponse
        }
    }
)


# --- GET --- #
@workspace_api_router.get(path="",
                          responses={
                              status.HTTP_200_OK: {
                                  'model': WorkspaceResponse
                              }
                          }
                          )
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

    return WorkspaceListResponse(data=workspaces, totalElements=count, page=data.page, size=data.size)


# --- PATCH --- #
@workspace_api_router.patch(path="",
                            status_code=status.HTTP_202_ACCEPTED)
async def patch_workspace(data: UpdateWorkspaceRequest, response: Response):
    try:
        update_workspace(data=data.model_dump(exclude_unset=True, exclude_none=True))
    except Exception as e:
        raise GenericApiException(status_code=status.HTTP_400_BAD_REQUEST, detail=str(e)) from e

    response.status_code = status.HTTP_202_ACCEPTED
    return None

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
