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
from meilisearch import Client
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

    SECURITY_JWT_SIGNING_KEY: str
    JWT_ALGORITHM: str

    URL: str
    SECURITY_CORS_ORIGINS: str
    API_URL: str
    IDP_URL: str
    WEBDAV_URL: str
    CONVERSION_URL: str
    LANGUAGE_URL: str
    MOSAIC_URL: str

    SEARCH_URL: str

    LOG_LEVEL: str = "INFO"
    LOG_FORMAT: str = "PLAIN"

settings = Settings()

conn = connect(conninfo=f"postgres://{settings.POSTGRES_USER}:"
                        f"{settings.POSTGRES_PASSWORD}@"
                        f"{settings.POSTGRES_URL}:"
                        f"{settings.POSTGRES_PORT}/"
                        f"{settings.POSTGRES_NAME}",
               row_factory=dict_row,
               autocommit=True
               )

meilisearch_client = Client(settings.SEARCH_URL)


def camel_to_snake(data: str):
    return reduce(lambda x, y: x + ('_' if y.isupper() else '') + y, data).lower()


def parse_sql_update_query(tablename: str, data: dict):
    date_fields = ", ".join(f'{camel_to_snake(k)} = \'{v.strftime("%Y-%m-%dT%H:%M:%SZ")}\'' for
                            k, v in data.items() if k != "id" and isinstance(v, datetime.datetime))
    data_fields = ", ".join(f'{camel_to_snake(k)} = \'{v}\'' for
                            k, v in data.items() if k != "id" and not isinstance(v, datetime.datetime))
    return (f'UPDATE "{tablename}" SET {data_fields} WHERE id = \'{data["id"]}\''
            f';') if date_fields == '' else (f'UPDATE "{tablename}" SET '
                                             f'{data_fields}, {date_fields} WHERE id = \'{data["id"]}\';')


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
                                           key=settings.SECURITY_JWT_SIGNING_KEY,
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
