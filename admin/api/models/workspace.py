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
class WorkspaceRequest(GenericRequest):
    pass


class WorkspaceListRequest(GenericPaginationRequest):
    pass


class OrganizationWorkspaceListRequest(WorkspaceRequest, WorkspaceListRequest):
    pass


# --- RESPONSE MODELS --- #
class WorkspaceResponse(GenericResponse):
    name: str
    organization_id: str
    storage_capacity: float
    root_id: str
    bucket: str
    create_time: datetime.datetime
    update_time: datetime.datetime


class WorkspaceListResponse(GenericListResponse):
    data: List[WorkspaceResponse]
