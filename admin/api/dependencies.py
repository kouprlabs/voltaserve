from typing import Optional

from pydantic_settings import BaseSettings


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

    class Config:
        env_file = "C:/Users/lobod/PycharmProjects/voltaserve/admin/api/.env"


settings = Settings()
