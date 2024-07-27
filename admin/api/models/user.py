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
class UserRequest(GenericRequest):
    pass


class UserListRequest(GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class UserResponse(GenericResponse):
    fullName: str
    username: str
    email: EmailStr
    isEmailConfirmed: bool
    picture: str | None
    createTime: datetime.datetime
    updateTime: datetime.datetime


class UserListResponse(GenericListResponse):
    data: List[UserResponse]
