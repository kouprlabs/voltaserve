from typing import Annotated

from fastapi import APIRouter, Depends, status

from ..database import fetch_group, fetch_groups
from ..exceptions import GenericNotFoundException
from ..models import GenericNotFoundResponse, GroupResponse, GroupListRequest, GroupListResponse, GroupRequest

group_api_router = APIRouter(
    prefix='/group',
    tags=['group'],
    responses={
        status.HTTP_404_NOT_FOUND: {
            'model': GenericNotFoundResponse
        }
    }
)


# --- GET --- #
@group_api_router.get(path="/",
                      responses={
                          status.HTTP_200_OK: {
                              'model': GroupResponse
                          }}
                      )
async def get_group(data: Annotated[GroupRequest, Depends()]):
    group = fetch_group(_id=data.id)
    if group is None:
        raise GenericNotFoundException(detail=f'Groups with id={data.id} does not exist')

    return GroupResponse(**group)


@group_api_router.get(path="/all",
                      responses={
                          status.HTTP_200_OK: {
                              'model': GroupListResponse
                          }
                      }
                      )
async def get_all_groups(data: Annotated[GroupListRequest, Depends()]):
    groups = fetch_groups(page=data.page, size=data.size)
    if groups is None:
        raise GenericNotFoundException(detail='This instance has no groups')

    return GroupListResponse(groups=groups)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
