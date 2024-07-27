# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

import time
from typing import Optional

import jwt
from fastapi import Request
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from pydantic import Field
from pydantic_settings import BaseSettings

from .exceptions import GenericForbiddenException


class Settings(BaseSettings):
    db_host: str
    db_port: int
    db_name: str
    db_user: str
    db_password: Optional[str]

    host: str
    workers: int
    port: int

    jwt_secret: str
    jwt_algorithm: str

    url: str
    cors_origins: str

    admin_token: Optional[str] = Field('')
    admin_token_expiration: int

    class Config:
        env_file = ".env"


settings = Settings()


class JWTBearer(HTTPBearer):
    def __init__(self, auto_error: bool = True):
        super(JWTBearer, self).__init__(auto_error=auto_error)

    async def __call__(self, request: Request):
        jwt.api_jws.PyJWS.header_typ = False
        credentials: HTTPAuthorizationCredentials = await super(JWTBearer, self).__call__(request)
        if credentials:
            if not credentials.scheme == "Bearer":
                raise GenericForbiddenException(detail="Invalid authentication scheme.")

            try:
                decoded_token = jwt.decode(jwt=credentials.credentials,
                                           key=settings.jwt_secret,
                                           algorithms=[settings.jwt_algorithm],
                                           audience=settings.url,
                                           issuer=settings.url,
                                           verify=True)

            except Exception as e:
                raise GenericForbiddenException(detail=str(e))

            if decoded_token['sub'] != 'SUPERUSER':
                raise GenericForbiddenException(detail='Invalid subject')
            if decoded_token['iat'] > time.time():
                raise GenericForbiddenException(detail='Token issued in the future')
            if decoded_token['exp'] < time.time():
                raise GenericForbiddenException(detail='Token expired')

            if credentials.credentials != settings.admin_token:
                raise GenericForbiddenException(detail="You are not a superuser :D")

            return credentials.credentials
        else:
            raise GenericForbiddenException(detail="Invalid token")
