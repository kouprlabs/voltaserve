from typing import Annotated

from fastapi import APIRouter, Depends, status

from admin.api.models.generic import GenericNotFoundResponse
from admin.api.models.user import UserListRequest, UserListResponse, UserRequest, UserResponse
from admin.api.database.user import fetch_user, fetch_users

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
@users_api_router.get(path="/",
                      responses={
                          status.HTTP_200_OK: {
                              'model': UserListResponse
                          }}
                      )
async def get_user(data: Annotated[UserRequest, Depends()]):
    user = fetch_user(_id=data.id)
    if user is None:
        return GenericNotFoundResponse(message=f'User with id={data.id} does not exist')

    return UserResponse(**user)


@users_api_router.get(path="/all",
                      responses={
                          status.HTTP_200_OK: {
                              'model': UserListResponse
                          }
                      }
                      )
async def get_all_users(data: Annotated[UserListRequest, Depends()]):
    users = fetch_users(page=data.page, size=data.size)
    if users is None:
        return GenericNotFoundResponse(message=f'This instance has no users')

    return UserListResponse(users=users)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
