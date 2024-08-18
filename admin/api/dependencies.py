# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

from functools import reduce
from typing import Optional

import jwt
from fastapi import Request
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
from psycopg import connect
from psycopg.rows import dict_row
from pydantic_settings import BaseSettings

from .errors import GenericForbiddenException


class Settings(BaseSettings):
    POSTGRES_URL: str
    POSTGRES_PORT: int
    POSTGRES_NAME: str
    POSTGRES_USER: str
    POSTGRES_PASSWORD: Optional[str]

    HOST: str
    WORKERS: int
    PORT: int

    JWT_SECRET: str
    JWT_ALGORITHM: str

    URL: str
    CORS_ORIGINS: str
    REDIS_URL: str
    MEILISEARCH_URL: str
    API_URL: str
    IDP_URL: str
    WEBDAV_URL: str
    CONVERSION_URL: str
    LANGUAGE_URL: str
    MOSAIC_URL: str
    MINIO_URL: str

    class Config:
        env_file = "C:/Users/lobod/PycharmProjects/voltaserve/admin/api/.env"


settings = Settings()

conn = connect(conninfo=f"postgres://{settings.POSTGRES_USER}:"
                        f"{settings.POSTGRES_PASSWORD}@"
                        f"{settings.POSTGRES_URL}:"
                        f"{settings.POSTGRES_PORT}/"
                        f"{settings.POSTGRES_NAME}",
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

            if not decoded_token['is_admin']:
                raise GenericForbiddenException(detail='Dude, u aint admin')

            return credentials.credentials
        else:
            raise GenericForbiddenException(detail="Invalid token")
