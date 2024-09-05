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

from pydantic import EmailStr, Field

from .generic import (
    GenericPaginationRequest,
    GenericResponse,
    GenericListResponse,
    GenericRequest,
)


# --- REQUEST MODELS --- #
class UserRequest(GenericRequest):
    pass


class UserListRequest(GenericPaginationRequest):
    pass


class UserOrganizationRequest(GenericRequest):
    pass


class UserOrganizationListRequest(UserOrganizationRequest, GenericPaginationRequest):
    pass


class UserWorkspaceRequest(GenericRequest):
    pass


class UserWorkspaceListRequest(UserWorkspaceRequest, GenericPaginationRequest):
    pass


class UserGroupRequest(GenericRequest):
    pass


class UserGroupListRequest(UserGroupRequest, GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class UserResponse(GenericResponse):
    fullName: str
    username: str
    email: EmailStr
    isEmailConfirmed: bool
    picture: str | None = Field(None)
    createTime: datetime.datetime
    updateTime: datetime.datetime


class UserListResponse(GenericListResponse):
    data: List[UserResponse]


class UserOrganizationResponse(GenericResponse):
    permission: str
    organizationId: str
    organizationName: str
    createTime: datetime.datetime


class UserWorkspaceResponse(GenericResponse):
    permission: str
    workspaceId: str
    workspaceName: str
    createTime: datetime.datetime


class UserGroupResponse(GenericResponse):
    permission: str
    groupId: str
    groupName: str
    createTime: datetime.datetime


class UserOrganizationListResponse(GenericListResponse):
    data: List[UserOrganizationResponse]


class UserWorkspaceListResponse(GenericListResponse):
    data: List[UserWorkspaceResponse]


class UserGroupListResponse(GenericListResponse):
    data: List[UserGroupResponse]
