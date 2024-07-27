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

from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class OrganizationUserRequest(GenericRequest):
    pass


class OrganizationUserListRequest(GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class OrganizationUserResponse(GenericResponse):
    organization_id: str
    user_id: str
    create_time: datetime.datetime


class OrganizationUserListResponse(GenericListResponse):
    data: List[OrganizationUserResponse]
