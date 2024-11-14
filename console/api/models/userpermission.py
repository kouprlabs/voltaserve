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

from pydantic import BaseModel

from .generic import (
    GenericPaginationRequest,
    GenericResponse,
    GenericListResponse,
    GenericRequest,
)


# --- REQUEST MODELS --- #
class UserPermissionRequest(GenericRequest):
    pass


class UserPermissionListRequest(GenericPaginationRequest):
    pass


class UserPermissionGrantRequest(BaseModel):
    user_id: str
    resource_id: str
    permission: str


class UserPermissionRevokeRequest(BaseModel):
    user_id: str
    resource_id: str


# --- RESPONSE MODELS --- #
class UserPermissionResponse(GenericResponse):
    user_id: str
    resource_id: str
    permission: str
    create_time: datetime.datetime


class UserPermissionListResponse(GenericListResponse):
    data: List[UserPermissionResponse]
