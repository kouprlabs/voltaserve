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

from ..database import fetch_workspace, fetch_workspace_count, fetch_workspaces
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
    GenericAcceptedResponse,
    GenericUnexpectedErrorResponse,
    WorkspaceListRequest,
    WorkspaceListResponse,
    WorkspaceRequest,
    WorkspaceResponse,
    WorkspaceSearchRequest,
)

workspace_api_router = APIRouter(
    prefix="/workspace",
    tags=["workspace"],
    responses={
        status.HTTP_500_INTERNAL_SERVER_ERROR: {"model": GenericUnexpectedErrorResponse},
        status.HTTP_202_ACCEPTED: {"model": GenericAcceptedResponse},
    },
    dependencies=[Depends(JWTBearer())],
)
logger = base_logger.getChild("workspace")


@workspace_api_router.get(path="", responses={status.HTTP_200_OK: {"model": WorkspaceResponse}})
async def get_workspace(data: Annotated[WorkspaceRequest, Depends()]):
    try:
        workspace = fetch_workspace(_id=data.id)
        return WorkspaceResponse(**workspace)
    except NotFoundException as e:
        logger.error(e)
        return NotFoundError(message=str(e))
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


@workspace_api_router.get(path="/count", responses={status.HTTP_200_OK: {"model": CountResponse}})
async def get_workspace_count():
    try:
        return CountResponse(**fetch_workspace_count())
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


@workspace_api_router.get(
    path="/all",
    responses={status.HTTP_200_OK: {"model": WorkspaceListResponse}},
    response_model=WorkspaceListResponse,
    response_model_exclude_none=True,
)
async def get_all_workspaces(data: Annotated[WorkspaceListRequest, Depends()], user_id: str = Depends(get_user_id)):
    try:
        workspaces, count = fetch_workspaces(user_id=user_id, page=data.page, size=data.size)
        return WorkspaceListResponse(
            data=workspaces,
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


@workspace_api_router.get(path="/search", responses={status.HTTP_200_OK: {"model": WorkspaceListResponse}})
async def get_search_workspaces(data: Annotated[WorkspaceSearchRequest, Depends()]):
    try:
        workspaces = meilisearch_client.index("workspace").search(
            data.query, {"page": data.page, "hitsPerPage": data.size}
        )
        hits = []
        for workspace in workspaces["hits"]:
            try:
                workspace = fetch_workspace(workspace["id"])
                hits.append(workspace)
            except NotFoundException:
                pass
        count = len(workspaces["hits"])
        return WorkspaceListResponse(
            data=hits,
            totalElements=count,
            totalPages=(count + data.size - 1) // data.size,
            page=data.page,
            size=data.size,
        )
    except NotFoundException as e:
        return NotFoundError(message=str(e))
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()
