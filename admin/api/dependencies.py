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
from functools import reduce
from typing import Optional

import jwt
from fastapi import Request
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from psycopg import connect
from psycopg.rows import dict_row
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

    idp_url: str
    api_url: str
    redis_url: str

    class Config:
        env_file = "C:/Users/lobod/PycharmProjects/voltaserve/admin/api/.env"


settings = Settings()

conn = connect(conninfo=f"postgres://{settings.db_user}:"
                        f"{settings.db_password}@"
                        f"{settings.db_host}:"
                        f"{settings.db_port}/"
                        f"{settings.db_name}",
               row_factory=dict_row,
               autocommit=True
               )


def camel_to_snake(data: str):
    return reduce(lambda x, y: x + ('_' if y.isupper() else '') + y, data).lower()


def parse_sql_update_query(tablename: str, data: dict):
    return (f'UPDATE "{tablename}" SET {", ".join(f'{camel_to_snake(k)} = '
                                                  f'\'{v}\'' for k, v in data.items() if k != "id")} '
            f'WHERE id = \'{data["id"]}\';')


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
                raise GenericForbiddenException(detail=str(e)) from e

            if decoded_token['iat'] > time.time():
                raise GenericForbiddenException(detail='Token issued in the future')
            if decoded_token['exp'] < time.time():
                raise GenericForbiddenException(detail='Token expired')

            return credentials.credentials
        else:
            raise GenericForbiddenException(detail="Invalid token")
