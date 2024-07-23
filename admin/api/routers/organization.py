from typing import Annotated

from fastapi import APIRouter, Depends, status

from admin.api.database.organization import fetch_organization, fetch_organizations
from admin.api.models.organization import OrganizationResponse, OrganizationRequest, OrganizationListResponse, \
    OrganizationListRequest

organization_api_router = APIRouter(
    prefix='/organization',
    tags=['organization']
)


# --- GET --- #
@organization_api_router.get(path="/",
                             responses={
                                 status.HTTP_200_OK: {
                                     'model': OrganizationResponse
                                 }}
                             )
async def get_user(data: Annotated[OrganizationRequest, Depends()]):
    return fetch_organization(_id=data.id)


@organization_api_router.get(path="/all",
                             responses={
                                 status.HTTP_200_OK: {
                                     'model': OrganizationListResponse
                                 }
                             }
                             )
async def get_all_users(data: Annotated[OrganizationListRequest, Depends()]):
    return fetch_organizations(page=data.page, size=data.size)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
