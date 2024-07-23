from typing import Annotated

from fastapi import APIRouter, Depends, status

from admin.api.database.invitation import fetch_invitation, fetch_invitations
from admin.api.models.invitation import InvitationResponse, InvitationListRequest, InvitationListResponse, InvitationRequest

invitation_api_router = APIRouter(
    prefix='/invitation',
    tags=['invitation']
)


# --- GET --- #
@invitation_api_router.get(path="/",
                           responses={
                               status.HTTP_200_OK: {
                                   'model': InvitationResponse
                               }}
                           )
async def get_invitation(data: Annotated[InvitationRequest, Depends()]):
    return fetch_invitation(_id=data.id)


@invitation_api_router.get(path="/all",
                           responses={
                               status.HTTP_200_OK: {
                                   'model': InvitationListResponse
                               }
                           }
                           )
async def get_all_invitations(data: Annotated[InvitationListRequest, Depends()]):
    return fetch_invitations(page=data.page, size=data.size)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
