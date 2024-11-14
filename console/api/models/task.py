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

from .generic import (
    GenericPaginationRequest,
    GenericResponse,
    GenericListResponse,
    GenericRequest,
)


class TaskRequest(GenericRequest):
    pass


class TaskListRequest(GenericPaginationRequest):
    pass


class TaskResponse(GenericResponse):
    name: str
    error: str
    percentage: float
    is_complete: bool
    is_indeterminate: bool
    user_id: str
    status: str
    payload: str
    task_id: str
    create_time: datetime.datetime
    update_time: datetime.datetime


class TaskListResponse(GenericListResponse):
    data: List[TaskResponse]
