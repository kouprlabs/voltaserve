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
class OrganizationRequest(GenericRequest):
    pass


class OrganizationListRequest(GenericPaginationRequest):
    pass


class UserOrganizationListRequest(OrganizationRequest, OrganizationListRequest):
    pass


# --- RESPONSE MODELS --- #
class OrganizationResponse(GenericResponse):
    name: str
    create_time: datetime.datetime
    update_time: datetime.datetime


class OrganizationListResponse(GenericListResponse):
    data: List[OrganizationResponse]


class UserOrganizationResponse(OrganizationResponse):
    permission: str


class UserOrganizationListResponse(GenericListResponse):
    data: List[UserOrganizationResponse]
