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

from pydantic import Field

from .generic import (
    GenericListResponse,
    GenericPaginationRequest,
    GenericRequest,
    GenericResponse,
    GenericSearchRequest,
)


class OrganizationRequest(GenericRequest):
    pass


class OrganizationListRequest(GenericPaginationRequest):
    pass


class OrganizationSearchRequest(GenericSearchRequest):
    pass


class OrganizationWorkspaceListRequest(GenericRequest, GenericPaginationRequest):
    pass


class OrganizationUserListRequest(GenericRequest, GenericPaginationRequest):
    pass


class OrganizationGroupListRequest(GenericRequest, GenericPaginationRequest):
    pass


class OrganizationResponse(GenericResponse):
    name: str
    createTime: datetime.datetime = Field(None)
    updateTime: datetime.datetime = Field(None)


class OrganizationUserResponse(GenericResponse):
    username: str
    picture: str | None = Field(None)
    permission: str
    createTime: datetime.datetime


class OrganizationWorkspaceResponse(GenericResponse):
    name: str
    createTime: datetime.datetime


class OrganizationGroupResponse(GenericResponse):
    name: str
    createTime: datetime.datetime


class OrganizationListResponse(GenericListResponse):
    data: List[OrganizationResponse]


class OrganizationUserListResponse(GenericListResponse):
    data: List[OrganizationUserResponse]


class OrganizationWorkspaceListResponse(GenericListResponse):
    data: List[OrganizationWorkspaceResponse]


class OrganizationGroupListResponse(GenericListResponse):
    data: List[OrganizationGroupResponse]
