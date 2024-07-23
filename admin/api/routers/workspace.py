from typing import Annotated

from fastapi import APIRouter, Depends, status, HTTPException

from admin.api.database.workspace import fetch_workspace, fetch_workspaces
from admin.api.models.generic import GenericNotFoundResponse
from admin.api.models.workspace import WorkspaceResponse, WorkspaceRequest, WorkspaceListResponse, WorkspaceListRequest

workspace_api_router = APIRouter(
    prefix='/workspace',
    tags=['workspace'],
    responses={
        status.HTTP_404_NOT_FOUND: {
            'model': GenericNotFoundResponse
        }
    }
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
        return GenericNotFoundResponse(message=f'Workspace with id={data.id} does not exist')

    return WorkspaceResponse(**workspace)


@workspace_api_router.get(path="/all",
                          responses={
                              status.HTTP_200_OK: {
                                  'model': WorkspaceListResponse
                              }
                          }
                          )
async def get_all_workspaces(data: Annotated[WorkspaceListRequest, Depends()]):
    workspaces = fetch_workspaces(page=data.page, size=data.size)
    if workspaces is None:
        return GenericNotFoundResponse(message=f'This instance has no workspaces')

    return WorkspaceListResponse(workspaces=workspaces)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
