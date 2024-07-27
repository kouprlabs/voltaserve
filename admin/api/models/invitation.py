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

from pydantic import EmailStr

from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class InvitationRequest(GenericRequest):
    pass


class InvitationListRequest(GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class InvitationResponse(GenericResponse):
    organization_id: str
    owner_id: str
    email: EmailStr
    status: str
    create_time: datetime.datetime
    update_time: datetime.datetime


class InvitationListResponse(GenericListResponse):
    data: List[InvitationResponse]
