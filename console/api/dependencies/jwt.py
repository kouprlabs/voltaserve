# Copyright (c) 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

import jwt
from fastapi import Request
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials

from . import settings
from ..errors import GenericForbiddenException


class JWTBearer(HTTPBearer):
    def __init__(self, auto_error: bool = True):
        super(JWTBearer, self).__init__(auto_error=auto_error)

    async def __call__(self, request: Request):
        jwt.api_jws.PyJWS.header_typ = False
        credentials: HTTPAuthorizationCredentials = await super(
            JWTBearer, self
        ).__call__(request)
        if credentials:
            if not credentials.scheme == "Bearer":
                raise GenericForbiddenException(detail="Invalid authentication scheme.")

            try:
                decoded_token = jwt.decode(
                    jwt=credentials.credentials,
                    key=settings.SECURITY_JWT_SIGNING_KEY,
                    algorithms=[settings.JWT_ALGORITHM],
                    audience=settings.URL,
                    issuer=settings.URL,
                    verify=True,
                )

            except Exception as e:
                raise GenericForbiddenException(detail=str(e)) from e

            if not decoded_token["is_admin"]:
                raise GenericForbiddenException(detail="User is not admin.")

            return credentials.credentials
        else:
            raise GenericForbiddenException(detail="Invalid token.")
