# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.
from typing import Optional

from pydantic_settings import BaseSettings


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


settings = Settings(_env_file="./api/.env")
