# Copyright (c) 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

from typing import List

from pydantic import BaseModel, AnyHttpUrl

from . import GenericListResponse, GenericRequest


class VersionRequest(GenericRequest):
    pass


class VersionResponse(BaseModel):
    name: str
    currentVersion: str
    latestVersion: str
    updateAvailable: bool | None
    location: AnyHttpUrl


class VersionsResponse(GenericListResponse):
    data: List[VersionResponse]
