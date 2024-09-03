# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

from typing import List

from pydantic import BaseModel, AnyHttpUrl

from . import GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class VersionRequest(GenericRequest):
    pass


# --- RESPONSE MODELS --- #
class VersionResponse(BaseModel):
    name: str
    currentVersion: str
    latestVersion: str
    updateAvailable: bool | None
    location: AnyHttpUrl


class VersionsResponse(GenericListResponse):
    data: List[VersionResponse]
