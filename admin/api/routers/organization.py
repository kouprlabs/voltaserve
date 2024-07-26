from typing import Annotated

from fastapi import APIRouter, Depends, status

from ..database import fetch_organization, fetch_organizations
from ..exceptions import GenericNotFoundException
from ..models import GenericNotFoundResponse, OrganizationResponse, OrganizationRequest, OrganizationListResponse, \
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
    organization = fetch_organization(organization_id=data.id)
    if organization is None:
        raise GenericNotFoundException(detail=f'Organization with id={data.id} does not exist')

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
    if organizations is None:
        raise GenericNotFoundException(detail='This instance has no organizations')

    return OrganizationListResponse(data=organizations, totalElements=count['count'], page=data.page, size=data.size)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
