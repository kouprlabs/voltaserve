from typing import Annotated

from fastapi import APIRouter, Depends, status

from admin.api.database.group import fetch_group, fetch_groups
from admin.api.models.group import GroupResponse, GroupListRequest, GroupListResponse, GroupRequest

group_api_router = APIRouter(
    prefix='/group',
    tags=['group']
)


# --- GET --- #
@group_api_router.get(path="/",
                      responses={
                          status.HTTP_200_OK: {
                              'model': GroupResponse
                          }}
                      )
async def get_user(data: Annotated[GroupRequest, Depends()]):
    return fetch_group(_id=data.id)


@group_api_router.get(path="/all",
                      responses={
                          status.HTTP_200_OK: {
                              'model': GroupListResponse
                          }
                      }
                      )
async def get_all_users(data: Annotated[GroupListRequest, Depends()]):
    return fetch_groups(page=data.page, size=data.size)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
