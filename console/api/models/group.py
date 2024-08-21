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

from ..models.organization import OrganizationResponse
from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class GroupRequest(GenericRequest):
    pass


class GroupListRequest(GenericPaginationRequest):
    pass


class UpdateGroupRequest(GenericRequest):
    name: str | None = Field(None)
    updateTime: datetime.datetime | None = Field(default_factory=datetime.datetime.now)


# --- RESPONSE MODELS --- #
class GroupResponse(GenericResponse):
    name: str
    organization: OrganizationResponse
    createTime: datetime.datetime = Field(None)
    updateTime: datetime.datetime = Field(None)


class GroupListResponse(GenericListResponse):
    data: List[GroupResponse]
