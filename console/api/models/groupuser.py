# Copyright (c) 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

import datetime
from typing import List

from .generic import (
    GenericPaginationRequest,
    GenericResponse,
    GenericListResponse,
    GenericRequest,
)


class UserGroupRequest(GenericRequest):
    pass


class UserGroupListRequest(GenericPaginationRequest):
    pass


class UserGroupResponse(GenericResponse):
    group_id: str
    user_id: str
    create_time: datetime.datetime


class UserGroupListResponse(GenericListResponse):
    data: List[UserGroupResponse]
