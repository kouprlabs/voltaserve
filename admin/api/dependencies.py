import time
from typing import Optional

import jwt
from fastapi import Request
from fastapi.security import HTTPBearer, HTTPAuthorizationCredentials
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

    class Config:
        env_file = "C:/Users/lobod/PycharmProjects/voltaserve/admin/api/.env"


settings = Settings()


class JWTBearer(HTTPBearer):
    def __init__(self, auto_error: bool = True):
        super(JWTBearer, self).__init__(auto_error=auto_error)

    async def __call__(self, request: Request):
        credentials: HTTPAuthorizationCredentials = await super(JWTBearer, self).__call__(request)
        if credentials:
            if not credentials.scheme == "Bearer":
                raise GenericForbiddenException(detail="Invalid authentication scheme.")
            try:
                decoded_token = jwt.decode(jwt=credentials.credentials,
                                           key=settings.jwt_secret,
                                           algorithms=[settings.jwt_algorithm])
                print(decoded_token)
            except Exception as e:
                raise GenericForbiddenException(detail=str(e))

            if decoded_token['iat'] > time.time():
                raise GenericForbiddenException(detail='Token issued in the future')
            if decoded_token['exp'] < time.time():
                raise GenericForbiddenException(detail='Token expired')
            if decoded_token['iss'] != settings.url != decoded_token['aud']:
                raise GenericForbiddenException(detail='Token issued for another instance')

            return credentials.credentials
        else:
            raise GenericForbiddenException(detail="Invalid token")
