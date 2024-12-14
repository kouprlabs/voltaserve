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
from ..database import (
    fetch_user_organizations,
    fetch_user_workspaces,
    fetch_user_groups,
    fetch_user_count,
)
from ..dependencies import JWTBearer
from ..log import base_logger
from ..errors import (
    NotFoundError,
    EmptyDataException,
    NoContentError,
    NotFoundException,
    UnknownApiError,
)
from ..models import (
    UserOrganizationListRequest,
    UserOrganizationListResponse,
    UserWorkspaceListResponse,
    UserWorkspaceListRequest,
    UserGroupListResponse,
    UserGroupListRequest,
    CountResponse,
)

users_api_router = APIRouter(
    prefix="/user", tags=["user"], dependencies=[Depends(JWTBearer())]
)
logger = base_logger.getChild("user")


@users_api_router.get(
    path="/organizations",
    responses={status.HTTP_200_OK: {"model": UserOrganizationListResponse}},
)
async def get_user_organizations(
    data: Annotated[UserOrganizationListRequest, Depends()]
):
    try:
        organizations, count = fetch_user_organizations(
            user_id=data.id, page=data.page, size=data.size
        )
        return UserOrganizationListResponse(
            data=organizations,
            totalElements=count,
            totalPages=(count + data.size - 1) // data.size,
            page=data.page,
            size=data.size,
        )
    except EmptyDataException as e:
        logger.error(e)
        return NoContentError()
    except NotFoundException as e:
        logger.error(e)
        return NotFoundError(message=str(e))
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


@users_api_router.get(
    path="/workspaces",
    responses={status.HTTP_200_OK: {"model": UserWorkspaceListResponse}},
)
async def get_user_workspaces(data: Annotated[UserWorkspaceListRequest, Depends()]):
    try:
        workspaces, count = fetch_user_workspaces(
            user_id=data.id, page=data.page, size=data.size
        )
        return UserWorkspaceListResponse(
            data=workspaces,
            totalElements=count,
            totalPages=(count + data.size - 1) // data.size,
            page=data.page,
            size=data.size,
        )
    except EmptyDataException as e:
        logger.error(e)
        return NoContentError()
    except NotFoundException as e:
        logger.error(e)
        return NotFoundError(message=str(e))
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


@users_api_router.get(
    path="/groups", responses={status.HTTP_200_OK: {"model": UserGroupListResponse}}
)
async def get_user_groups(data: Annotated[UserGroupListRequest, Depends()]):
    try:
        groups, count = fetch_user_groups(
            user_id=data.id, page=data.page, size=data.size
        )
        return UserGroupListResponse(
            data=groups,
            totalElements=count,
            totalPages=(count + data.size - 1) // data.size,
            page=data.page,
            size=data.size,
        )
    except EmptyDataException as e:
        logger.error(e)
        return NoContentError()
    except NotFoundException as e:
        logger.error(e)
        return NotFoundError(message=str(e))
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


@users_api_router.get(
    path="/count", responses={status.HTTP_200_OK: {"model": CountResponse}}
)
async def get_user_count():
    try:
        return CountResponse(**fetch_user_count())
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()
