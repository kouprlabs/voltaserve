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

from pydantic import Field

from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class OrganizationRequest(GenericRequest):
    pass


class OrganizationListRequest(GenericPaginationRequest):
    pass


class UpdateOrganizationRequest(GenericRequest):
    name: str | None = Field(None)
    updateTime: datetime.datetime | None = Field(default_factory=datetime.datetime.now)


# --- RESPONSE MODELS --- #
class OrganizationResponse(GenericResponse):
    name: str
    createTime: datetime.datetime = Field(None)
    updateTime: datetime.datetime = Field(None)


class OrganizationListResponse(GenericListResponse):
    data: List[OrganizationResponse]
