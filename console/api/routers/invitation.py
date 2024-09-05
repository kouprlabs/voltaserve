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

from ..database import fetch_invitation, fetch_invitations, update_invitation
from ..dependencies import JWTBearer
from ..log import base_logger
from ..errors import (
    NotFoundError,
    NoContentError,
    EmptyDataException,
    NotFoundException,
    UnknownApiError,
)
from ..models import (
    InvitationResponse,
    InvitationListRequest,
    InvitationListResponse,
    InvitationRequest,
    ConfirmInvitationRequest,
    GenericAcceptedResponse,
)

invitation_api_router = APIRouter(
    prefix="/invitation", tags=["invitation"], dependencies=[Depends(JWTBearer())]
)


logger = base_logger.getChild("invitation")


# --- GET --- #
@invitation_api_router.get(
    path="", responses={status.HTTP_200_OK: {"model": InvitationResponse}}
)
async def get_invitation(data: Annotated[InvitationRequest, Depends()]):
    try:
        invitation = fetch_invitation(_id=data.id)

        return InvitationResponse(**invitation)
    except NotFoundException as e:
        logger.error(e)
        return NotFoundError(message=str(e))
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


@invitation_api_router.get(
    path="/all", responses={status.HTTP_200_OK: {"model": InvitationListResponse}}
)
async def get_all_invitations(data: Annotated[InvitationListRequest, Depends()]):
    try:
        invitations, count = fetch_invitations(page=data.page, size=data.size)

        return InvitationListResponse(
            data=invitations, totalElements=count, page=data.page, size=data.size
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


# --- PATCH --- #
@invitation_api_router.patch(
    path="",
    responses={status.HTTP_202_ACCEPTED: {"model": GenericAcceptedResponse}},
    status_code=status.HTTP_202_ACCEPTED,
)
async def patch_invitation(data: ConfirmInvitationRequest, response: Response):
    try:
        update_invitation(data.model_dump(exclude_unset=True, exclude_none=True))

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
