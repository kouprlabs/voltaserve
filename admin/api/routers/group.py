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

from ..database import fetch_groups, fetch_group
from ..database.group import update_group, fetch_group_count
from ..dependencies import JWTBearer
from ..errors import NotFoundError, NoContentError, EmptyDataException, NotFoundException, \
    UnknownApiError
from ..models import GroupResponse, GroupListRequest, GroupListResponse, GroupRequest, \
    UpdateGroupRequest, CountResponse

group_api_router = APIRouter(
    prefix='/group',
    tags=['group'],
    dependencies=[Depends(JWTBearer())]
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
    try:
        group = fetch_group(_id=data.id)

        return GroupResponse(**group)
    except NotFoundException as e:
        return NotFoundError(message=str(e))
    except Exception as e:
        return UnknownApiError(message=str(e))


@group_api_router.get(path="/count",
                      responses={
                          status.HTTP_200_OK: {
                              'model': CountResponse
                          }
                      }
                      )
async def get_group_count():
    try:
        return CountResponse(**fetch_group_count())
    except EmptyDataException:
        return NoContentError()
    except NotFoundException as e:
        return NotFoundError(message=str(e))
    except Exception as e:
        return UnknownApiError(message=str(e))


@group_api_router.get(path="/all",
                      responses={
                          status.HTTP_200_OK: {
                              'model': GroupListResponse
                          }
                      }
                      )
async def get_all_groups(data: Annotated[GroupListRequest, Depends()]):
    try:
        groups, count = fetch_groups(page=data.page, size=data.size)

        return GroupListResponse(data=groups, totalElements=count, page=data.page, size=data.size)
    except EmptyDataException:
        return NoContentError()
    except NotFoundException as e:
        return NotFoundError(message=str(e))
    except Exception as e:
        return UnknownApiError(message=str(e))


# --- PATCH --- #
@group_api_router.patch(path="",
                        status_code=status.HTTP_202_ACCEPTED)
async def patch_group(data: UpdateGroupRequest, response: Response):
    try:
        update_group(data=data.model_dump(exclude_unset=True, exclude_none=True))
    except NotFoundException as e:
        return NotFoundError(message=str(e))
    except Exception as e:
        return UnknownApiError(message=str(e))

    response.status_code = status.HTTP_202_ACCEPTED
    return None

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
