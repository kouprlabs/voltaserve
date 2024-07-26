from fastapi import HTTPException, status


class GenericNotFoundException(HTTPException):
    def __init__(self, detail='Not found'):
        super().__init__(status_code=status.HTTP_404_NOT_FOUND, detail=detail)


class GenericUnauthorizedException(HTTPException):
    def __init__(self, detail='Unauthorized'):
        super().__init__(status_code=status.HTTP_401_UNAUTHORIZED, detail=detail)


class GenericForbiddenException(HTTPException):
    def __init__(self, detail='Forbidden'):
        super().__init__(status_code=status.HTTP_403_FORBIDDEN, detail=detail)
