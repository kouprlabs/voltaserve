from typing import Annotated

from fastapi import APIRouter, Depends, status

from admin.api.database.organization import fetch_organization, fetch_organizations
from admin.api.models.generic import GenericNotFoundResponse
from admin.api.models.organization import OrganizationResponse, OrganizationRequest, OrganizationListResponse, \
    OrganizationListRequest

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
@organization_api_router.get(path="/",
                             responses={
                                 status.HTTP_200_OK: {
                                     'model': OrganizationResponse
                                 }}
                             )
async def get_organization(data: Annotated[OrganizationRequest, Depends()]):
    organization = fetch_organization(_id=data.id)
    if organization is None:
        return GenericNotFoundResponse(message=f'Organization with id={data.id} does not exist')

    return OrganizationResponse(**organization)


@organization_api_router.get(path="/all",
                             responses={
                                 status.HTTP_200_OK: {
                                     'model': OrganizationListResponse
                                 }
                             }
                             )
async def get_all_organizations(data: Annotated[OrganizationListRequest, Depends()]):
    organizations = fetch_organizations(page=data.page, size=data.size)
    if organizations is None:
        return GenericNotFoundResponse(message=f'This instance has no organizations')

    return OrganizationListResponse(organizations=organizations)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
