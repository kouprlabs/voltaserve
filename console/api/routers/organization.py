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
    fetch_organization,
    fetch_organization_count,
    fetch_organization_groups,
    fetch_organization_users,
    fetch_organization_workspaces,
    fetch_organizations,
)
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
    OrganizationGroupListRequest,
    OrganizationGroupListResponse,
    OrganizationListRequest,
    OrganizationListResponse,
    OrganizationRequest,
    OrganizationResponse,
    OrganizationSearchRequest,
    OrganizationUserListRequest,
    OrganizationUserListResponse,
    OrganizationWorkspaceListRequest,
    OrganizationWorkspaceListResponse,
)

organization_api_router = APIRouter(prefix="/organization", tags=["organization"], dependencies=[Depends(JWTBearer())])
logger = base_logger.getChild("organization")


@organization_api_router.get(path="", responses={status.HTTP_200_OK: {"model": OrganizationResponse}})
async def get_organization(data: Annotated[OrganizationRequest, Depends()], user_id: str = Depends(get_user_id)):
    try:
        organization = fetch_organization(organization_id=data.id, user_id=user_id)
        return OrganizationResponse(**organization)
    except NotFoundException as e:
        logger.error(e)
        return NotFoundError(message=str(e))
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


@organization_api_router.get(path="/count", responses={status.HTTP_200_OK: {"model": CountResponse}})
async def get_organization_count():
    try:
        return CountResponse(**fetch_organization_count())
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


@organization_api_router.get(
    path="/all",
    responses={status.HTTP_200_OK: {"model": OrganizationListResponse}},
    response_model=OrganizationListResponse,
    response_model_exclude_none=True,
)
async def get_all_organizations(
    data: Annotated[OrganizationListRequest, Depends()], user_id: str = Depends(get_user_id)
):
    try:
        organizations, count = fetch_organizations(user_id=user_id, page=data.page, size=data.size)
        return OrganizationListResponse(
            data=organizations,
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


@organization_api_router.get(path="/search", responses={status.HTTP_200_OK: {"model": OrganizationListResponse}})
async def get_search_organizations(
    data: Annotated[OrganizationSearchRequest, Depends()], user_id: str = Depends(get_user_id)
):
    try:
        organizations = meilisearch_client.index("organization").search(
            data.query, {"page": data.page, "hitsPerPage": data.size}
        )
        hits = []
        for organization in organizations["hits"]:
            try:
                organization = fetch_organization(organization["id"], user_id)
                hits.append(organization)
            except NotFoundException:
                pass
        count = len(organizations["hits"])
        return OrganizationListResponse(
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


@organization_api_router.get(
    path="/users",
    responses={status.HTTP_200_OK: {"model": OrganizationUserListResponse}},
)
async def get_organization_users(data: Annotated[OrganizationUserListRequest, Depends()]):
    try:
        users, count = fetch_organization_users(organization_id=data.id, page=data.page, size=data.size)
        return OrganizationUserListResponse(
            data=users,
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


@organization_api_router.get(
    path="/workspaces",
    responses={status.HTTP_200_OK: {"model": OrganizationWorkspaceListResponse}},
)
async def get_organization_workspaces(data: Annotated[OrganizationWorkspaceListRequest, Depends()]):
    try:
        workspaces, count = fetch_organization_workspaces(organization_id=data.id, page=data.page, size=data.size)
        return OrganizationWorkspaceListResponse(
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


@organization_api_router.get(
    path="/groups",
    responses={status.HTTP_200_OK: {"model": OrganizationGroupListResponse}},
)
async def get_organization_groups(data: Annotated[OrganizationGroupListRequest, Depends()]):
    try:
        groups, count = fetch_organization_groups(organization_id=data.id, page=data.page, size=data.size)
        return OrganizationWorkspaceListResponse(
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
