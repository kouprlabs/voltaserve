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
from ..dependencies import JWTBearer, meilisearch_client
from ..log import base_logger
from ..errors import NotFoundError, NoContentError, EmptyDataException, NotFoundException, \
    UnknownApiError
from ..models import GroupResponse, GroupListRequest, GroupListResponse, GroupRequest, \
    UpdateGroupRequest, CountResponse, GroupSearchRequest

group_api_router = APIRouter(
    prefix='/group',
    tags=['group'],
    dependencies=[Depends(JWTBearer())]
)

logger = base_logger.getChild("group")


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
        logger.error(e)
        return NotFoundError(message=str(e))
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


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
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


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
    except EmptyDataException as e:
        logger.error(e)
        return NoContentError()
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


@group_api_router.get(path="/search",
                      responses={
                          status.HTTP_200_OK: {
                              'model': GroupListResponse
                          }
                      }
                      )
async def get_search_groups(data: Annotated[GroupSearchRequest, Depends()]):
    try:
        groups = meilisearch_client.index('group').search(data.query, {'page': data.page, 'hitsPerPage': data.size})

        return GroupListResponse(data=(fetch_group(group['id']) for group in groups['hits']),
                                 totalElements=len(groups['hits']), page=data.page, size=data.size)
    except NotFoundException as e:
        logger.error(e)
        return NotFoundError(message=str(e))
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


# --- PATCH --- #
@group_api_router.patch(path="",
                        status_code=status.HTTP_202_ACCEPTED)
async def patch_group(data: UpdateGroupRequest, response: Response):
    try:
        update_group(data=data.model_dump(exclude_none=True))
        meilisearch_client.index('group').update_documents([{
            'id': data.id,
            'name': data.name,
            'updateTime': data.updateTime.strftime("%Y-%m-%dT%H:%M:%SZ")
        }])
    except NotFoundException as e:
        logger.error(e)
        return NotFoundError(message=str(e))
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()

    response.status_code = status.HTTP_202_ACCEPTED
    return None

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
