# Copyright (c) 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

from .generic import GenericTokenPayload, GenericTokenResponse


class TokenPayload(GenericTokenPayload):
    pass


# class TokenRequest(GenericTokenRequest):
#     Authorization: Header


class TokenResponse(GenericTokenResponse):
    authorized: bool = False
