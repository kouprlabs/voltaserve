# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.
import logging
from typing import Annotated

from fastapi import APIRouter, Depends, status, Response

from ..database import fetch_groups, fetch_group
from ..database.group import update_group
from ..exceptions import GenericNotFoundException, GenericApiException
from ..models import GenericNotFoundResponse, GroupResponse, GroupListRequest, GroupListResponse, GroupRequest, \
    UpdateGroupRequest

group_api_router = APIRouter(
    prefix='/group',
    tags=['group'],
    responses={
        status.HTTP_404_NOT_FOUND: {
            'model': GenericNotFoundResponse
        }
    },
)


# --- GET --- #
@group_api_router.get(path="",
                      responses={
                          status.HTTP_200_OK: {
                              'model': GroupResponse
                          }
                      }
                      )
async def get_group(data: Annotated[GroupRequest, Depends()]):
    group = fetch_group(_id=data.id)
    if group is None:
        raise GenericNotFoundException(detail=f'Group with id={data.id} does not exist')

    return GroupResponse(**group)


@group_api_router.get(path="/all",
                      responses={
                          status.HTTP_200_OK: {
                              'model': GroupListResponse
                          }
                      }
                      )
async def get_all_groups(data: Annotated[GroupListRequest, Depends()]):
    groups, count = fetch_groups(page=data.page, size=data.size)
    if groups is None or count is None:
        raise GenericNotFoundException(detail='This instance has no groups')

    return GroupListResponse(data=groups, totalElements=count, page=data.page, size=data.size)


# --- PATCH --- #
@group_api_router.patch(path="",
                        status_code=status.HTTP_202_ACCEPTED)
async def patch_group(data: UpdateGroupRequest, response: Response):
    try:
        update_group(data=data.model_dump(exclude_unset=True, exclude_none=True))
    except Exception as e:
        raise GenericApiException(status_code=status.HTTP_400_BAD_REQUEST, detail=str(e)) from e

    response.status_code = status.HTTP_202_ACCEPTED
    return None

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
