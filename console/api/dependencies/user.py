from fastapi import Depends

from . import JWTBearer

jwt_bearer = JWTBearer()


async def get_user_id(user_id: dict = Depends(jwt_bearer)):
    return user_id
