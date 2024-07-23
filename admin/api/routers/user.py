from typing import Annotated

from fastapi import APIRouter, Depends, status
from admin.api.models.user import UserListRequest, UserListResponse, UserRequest
from admin.api.database.user import fetch_user, fetch_users

users_api_router = APIRouter(
    prefix='/users'
)


# --- GET --- #
@users_api_router.get(path="/",
                      responses={
                          status.HTTP_200_OK: {
                              'model': UserListResponse
                          }}
                      )
async def get_user(data: Annotated[UserRequest, Depends()]):
    return fetch_user(_id=data.id)


@users_api_router.get(path="/all",
                      responses={
                          status.HTTP_200_OK: {
                              'model': UserListResponse
                          }
                      }
                      )
async def get_all_users(data: Annotated[UserListRequest, Depends()]):
    return fetch_users(page=data.page, size=data.size)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
