# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

import datetime
from typing import List

from pydantic import EmailStr, Field

from ..models.organization import OrganizationResponse
from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class InvitationRequest(GenericRequest):
    pass


class InvitationListRequest(GenericPaginationRequest):
    pass


class ConfirmInvitationRequest(GenericRequest):
    accept: bool
    updateTime: datetime.datetime | None = Field(default_factory=datetime.datetime.now)


# --- RESPONSE MODELS --- #
class InvitationResponse(GenericResponse):
    organization: OrganizationResponse
    ownerId: str
    email: EmailStr
    status: str
    createTime: datetime.datetime
    updateTime: datetime.datetime


class InvitationListResponse(GenericListResponse):
    data: List[InvitationResponse]
