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

from ..database import fetch_organization, fetch_organizations, fetch_organization_users, fetch_organization_workspaces
from ..database.organization import update_organization
from ..exceptions import GenericNotFoundException, GenericApiException
from ..models import GenericNotFoundResponse, OrganizationResponse, OrganizationRequest, OrganizationListResponse, \
    OrganizationListRequest, OrganizationUserListRequest, OrganizationUserListResponse, \
    WorkspaceListResponse, OrganizationWorkspaceListRequest, UpdateOrganizationRequest

organization_api_router = APIRouter(
    prefix='/organization',
    tags=['organization'],
    responses={
        status.HTTP_404_NOT_FOUND: {
            'model': GenericNotFoundResponse
        }
    }
)


# --- GET --- #
@organization_api_router.get(path="",
                             responses={
                                 status.HTTP_200_OK: {
                                     'model': OrganizationResponse
                                 }
                             }
                             )
async def get_organization(data: Annotated[OrganizationRequest, Depends()]):
    organization = fetch_organization(organization_id=data.id)
    if organization is None or organization == {}:
        raise GenericNotFoundException(detail=f'Organization with id={data.id} does not exist!')

    return OrganizationResponse(**organization)


@organization_api_router.get(path="/all",
                             responses={
                                 status.HTTP_200_OK: {
                                     'model': OrganizationListResponse
                                 }
                             }
                             )
async def get_all_organizations(data: Annotated[OrganizationListRequest, Depends()]):
    organizations, count = fetch_organizations(page=data.page, size=data.size)
    if organizations is None or organizations == {} or count == 0:
        raise GenericNotFoundException(detail='This instance has no organizations')

    return OrganizationListResponse(data=organizations, totalElements=count, page=data.page, size=data.size)


@organization_api_router.get(path="/users",
                             responses={
                                 status.HTTP_200_OK: {
                                     'model': OrganizationListResponse
                                 }
                             }
                             )
async def get_organization_users(data: Annotated[OrganizationUserListRequest, Depends()]):
    users, count = fetch_organization_users(organization_id=data.id, page=data.page, size=data.size)
    if users is None or users == {}:
        raise GenericNotFoundException(detail='This organization has no users')

    return OrganizationUserListResponse(data=users, totalElements=count, page=data.page, size=data.size)


@organization_api_router.get(path="/workspaces",
                             responses={
                                 status.HTTP_200_OK: {
                                     'model': WorkspaceListResponse
                                 }
                             }
                             )
async def get_organization_workspaces(data: Annotated[OrganizationWorkspaceListRequest, Depends()]):
    workspaces, count = fetch_organization_workspaces(organization_id=data.id, page=data.page, size=data.size)
    if workspaces is None or workspaces == {} or count == 0:
        raise GenericNotFoundException(detail='This organization has no workspaces')

    return WorkspaceListResponse(data=workspaces, totalElements=count, page=data.page, size=data.size)


# --- PATCH --- #
@organization_api_router.patch(path="",
                               status_code=status.HTTP_202_ACCEPTED)
async def patch_workspace(data: UpdateOrganizationRequest, response: Response):
    try:
        update_organization(data=data.model_dump(exclude_unset=True, exclude_none=True))
    except Exception as e:
        raise GenericApiException(status_code=status.HTTP_400_BAD_REQUEST, detail=str(e)) from e

    response.status_code = status.HTTP_202_ACCEPTED
    return None

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
