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

from .generic import GenericPaginationRequest, GenericListResponse, GenericRequest, BaseModel


# --- REQUEST MODELS --- #
class IndexRequest(GenericRequest):
    pass


class IndexListRequest(GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class IndexResponse(BaseModel):
    tablename: str
    indexname: str
    indexdef: str


class IndexListResponse(GenericListResponse):
    data: List[IndexResponse]
