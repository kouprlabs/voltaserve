from typing import Annotated

from fastapi import APIRouter, Depends, status

from admin.api.database.workspace import fetch_workspace, fetch_workspaces
from admin.api.models.workspace import WorkspaceResponse, WorkspaceRequest, WorkspaceListResponse, WorkspaceListRequest

workspace_api_router = APIRouter(
    prefix='/workspace'
)


# --- GET --- #
@workspace_api_router.get(path="/",
                          responses={
                              status.HTTP_200_OK: {
                                  'model': WorkspaceResponse
                              }}
                          )
async def get_workspace(data: Annotated[WorkspaceRequest, Depends()]):
    return fetch_workspace(_id=data.id)


@workspace_api_router.get(path="/all",
                          responses={
                              status.HTTP_200_OK: {
                                  'model': WorkspaceListResponse
                              }
                          }
                          )
async def get_all_workspaces(data: Annotated[WorkspaceListRequest, Depends()]):
    return fetch_workspaces(page=data.page, size=data.size)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
