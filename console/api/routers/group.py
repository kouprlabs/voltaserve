# Copyright (c) 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

from typing import Annotated

from fastapi import APIRouter, Depends, status

from ..database import fetch_group, fetch_groups
from ..database.group import fetch_group_count
from ..dependencies import JWTBearer, meilisearch_client
from ..dependencies.user import get_user_id
from ..errors import (
    EmptyDataException,
    NoContentError,
    NotFoundError,
    NotFoundException,
    UnknownApiError,
)
from ..log import base_logger
from ..models import (
    CountResponse,
    GroupListRequest,
    GroupListResponse,
    GroupRequest,
    GroupResponse,
    GroupSearchRequest,
)

group_api_router = APIRouter(prefix="/group", tags=["group"], dependencies=[Depends(JWTBearer())])
logger = base_logger.getChild("group")


@group_api_router.get(path="", responses={status.HTTP_200_OK: {"model": GroupResponse}})
async def get_group(data: Annotated[GroupRequest, Depends()], user_id: str = Depends(get_user_id)):
    try:
        group = fetch_group(_id=data.id, user_id=user_id)
        return GroupResponse(**group)
    except NotFoundException as e:
        logger.error(e)
        return NotFoundError(message=str(e))
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


@group_api_router.get(path="/count", responses={status.HTTP_200_OK: {"model": CountResponse}})
async def get_group_count():
    try:
        return CountResponse(**fetch_group_count())
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


@group_api_router.get(
    path="/all",
    responses={status.HTTP_200_OK: {"model": GroupListResponse}},
    response_model=GroupListResponse,
    response_model_exclude_none=True,
)
async def get_all_groups(data: Annotated[GroupListRequest, Depends()], user_id: str = Depends(get_user_id)):
    try:
        groups, count = fetch_groups(user_id=user_id, page=data.page, size=data.size)
        return GroupListResponse(
            data=groups,
            totalElements=count,
            totalPages=(count + data.size - 1) // data.size,
            page=data.page,
            size=data.size,
        )
    except EmptyDataException as e:
        logger.error(e)
        return NoContentError()
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


@group_api_router.get(path="/search", responses={status.HTTP_200_OK: {"model": GroupListResponse}})
async def get_search_groups(data: Annotated[GroupSearchRequest, Depends()], user_id: str = Depends(get_user_id)):
    try:
        groups = meilisearch_client.index("group").search(data.query, {"page": data.page, "hitsPerPage": data.size})
        hits = []
        for group in groups["hits"]:
            try:
                group = fetch_group(group["id"], user_id)
                hits.append(group)
            except NotFoundException:
                pass
        count = len(groups["hits"])
        return GroupListResponse(
            data=hits,
            totalElements=count,
            totalPages=(count + data.size - 1) // data.size,
            page=data.page,
            size=data.size,
        )
    except NotFoundException as e:
        logger.error(e)
        return NotFoundError(message=str(e))
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()
