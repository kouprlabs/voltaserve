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

from api.routers.exceptions import GenericForbiddenException


class Settings(BaseSettings):
    DB_HOST: str
    DB_PORT: int
    DB_NAME: str
    DB_USER: str
    DB_PASSWORD: Optional[str]

    HOST: str
    WORKERS: int
    PORT: int

    JWT_SECRET: str
    JWT_ALGORITHM: str

    URL: str
    CORS_ORIGINS: str

    IDP_URL: str
    API_URL: str
    REDIS_URL: str

    class Config:
        env_file = "C:/Users/lobod/PycharmProjects/voltaserve/admin/api/.env"


settings = Settings()

conn = connect(conninfo=f"postgres://{settings.DB_USER}:"
                        f"{settings.DB_PASSWORD}@"
                        f"{settings.DB_HOST}:"
                        f"{settings.DB_PORT}/"
                        f"{settings.DB_NAME}",
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
                                           key=settings.JWT_SECRET,
                                           algorithms=[settings.JWT_ALGORITHM],
                                           audience=settings.URL,
                                           issuer=settings.URL,
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
