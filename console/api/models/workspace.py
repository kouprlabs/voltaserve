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
from .generic import (
    GenericPaginationRequest,
    GenericResponse,
    GenericListResponse,
    GenericRequest,
    GenericSearchRequest,
)


# --- REQUEST MODELS --- #
class WorkspaceRequest(GenericRequest):
    pass


class UpdateWorkspaceRequest(GenericRequest):
    name: str | None = Field(None)
    organizationId: str | None = Field(None)
    storageCapacity: float | None = Field(None)
    rootId: str | None = Field(None)
    bucket: str | None = Field(None)
    updateTime: datetime.datetime | None = Field(default_factory=datetime.datetime.now)


class WorkspaceListRequest(GenericPaginationRequest):
    pass


class WorkspaceSearchRequest(GenericSearchRequest):
    pass


class WorkspaceOrganizationListRequest(WorkspaceRequest, WorkspaceListRequest):
    pass


# --- RESPONSE MODELS --- #
class WorkspaceResponse(GenericResponse):
    name: str
    organization: OrganizationResponse
    storageCapacity: float
    rootId: str | None = Field(None)
    bucket: str | None = Field(None)
    createTime: datetime.datetime
    updateTime: datetime.datetime


class WorkspaceListResponse(GenericListResponse):
    data: List[WorkspaceResponse]
