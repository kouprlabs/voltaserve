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

from ..database import (
    fetch_workspace,
    fetch_workspaces,
    update_workspace,
    fetch_workspace_count,
)
from ..dependencies import JWTBearer, meilisearch_client, redis_conn
from ..log import base_logger
from ..errors import (
    NotFoundError,
    EmptyDataException,
    NoContentError,
    NotFoundException,
    UnknownApiError,
)
from ..models import (
    WorkspaceResponse,
    WorkspaceRequest,
    WorkspaceListResponse,
    WorkspaceListRequest,
    UpdateWorkspaceRequest,
    GenericUnexpectedErrorResponse,
    GenericAcceptedResponse,
    CountResponse,
    WorkspaceSearchRequest,
)

workspace_api_router = APIRouter(
    prefix="/workspace",
    tags=["workspace"],
    responses={
        status.HTTP_500_INTERNAL_SERVER_ERROR: {
            "model": GenericUnexpectedErrorResponse
        },
        status.HTTP_202_ACCEPTED: {"model": GenericAcceptedResponse},
    },
    dependencies=[Depends(JWTBearer())],
)


logger = base_logger.getChild("workspace")


# --- GET --- #
@workspace_api_router.get(
    path="", responses={status.HTTP_200_OK: {"model": WorkspaceResponse}}
)
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


@workspace_api_router.get(
    path="/count", responses={status.HTTP_200_OK: {"model": CountResponse}}
)
async def get_workspace_count():
    try:
        return CountResponse(**fetch_workspace_count())
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


@workspace_api_router.get(
    path="/all", responses={status.HTTP_200_OK: {"model": WorkspaceListResponse}}
)
async def get_all_workspaces(data: Annotated[WorkspaceListRequest, Depends()]):
    try:
        workspaces, count = fetch_workspaces(page=data.page, size=data.size)

        return WorkspaceListResponse(
            data=workspaces, totalElements=count, page=data.page, size=data.size
        )
    except EmptyDataException as e:
        logger.error(e)
        return NoContentError()
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


@workspace_api_router.get(
    path="/search", responses={status.HTTP_200_OK: {"model": WorkspaceListResponse}}
)
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
        return WorkspaceListResponse(
            data=hits,
            totalElements=len(workspaces["hits"]),
            page=data.page,
            size=data.size,
        )
    except NotFoundException as e:
        return NotFoundError(message=str(e))
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


# --- PATCH --- #
@workspace_api_router.patch(path="", status_code=status.HTTP_202_ACCEPTED)
async def patch_workspace(data: UpdateWorkspaceRequest, response: Response):
    try:
        await redis_conn.delete(f"workspace:{data.id}")
        update_workspace(data=data.model_dump(exclude_none=True))
        meilisearch_client.index("workspace").update_documents(
            [
                {
                    "id": data.id,
                    "name": data.name,
                    "storageCapacity": data.storageCapacity,
                    "updateTime": data.updateTime.strftime("%Y-%m-%dT%H:%M:%SZ"),
                }
            ]
        )
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
