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

from fastapi import APIRouter, Depends, status

from ..database import fetch_invitation, fetch_invitations
from ..dependencies import JWTBearer
from ..exceptions import GenericNotFoundException
from ..models import GenericNotFoundResponse, InvitationResponse, InvitationListRequest, \
    InvitationListResponse, InvitationRequest

invitation_api_router = APIRouter(
    prefix='/invitation',
    tags=['invitation'],
    responses={
        status.HTTP_404_NOT_FOUND: {
            'model': GenericNotFoundResponse
        }
    },
    dependencies=[Depends(JWTBearer())]
)


# --- GET --- #
@invitation_api_router.get(path="/",
                           responses={
                               status.HTTP_200_OK: {
                                   'model': InvitationResponse
                               }}
                           )
async def get_invitation(data: Annotated[InvitationRequest, Depends()]):
    invitation = fetch_invitation(_id=data.id)(_id=data.id)
    if invitation is None:
        raise GenericNotFoundException(detail=f'Invitation with id={data.id} does not exist')

    return InvitationResponse(**invitation)


@invitation_api_router.get(path="/all",
                           responses={
                               status.HTTP_200_OK: {
                                   'model': InvitationListResponse
                               }
                           }
                           )
async def get_all_invitations(data: Annotated[InvitationListRequest, Depends()]):
    invitations, count = fetch_invitations(page=data.page, size=data.size)
    if invitations is None:
        raise GenericNotFoundException(detail='This instance has no invitations')

    return InvitationListResponse(data=invitations, totalElements=count['count'], page=data.page, size=data.size)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
