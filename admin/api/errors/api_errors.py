# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

from fastapi.responses import JSONResponse, Response
from fastapi import status
from .error_codes import errors


class GenericError(JSONResponse):
    def __init__(self,
                 status_code: int,
                 message: str,
                 user_message: str | None = None,
                 code: str | None = None):
        super().__init__(
            status_code=status_code,
            # headers=headers,
            content={
                'code': errors[status_code] if code is None else code,
                'status': status_code,
                'message': message,
                'userMessage': message if user_message is None else user_message,
                'moreInfo': f'https://voltaserve.com/docs/console/errors/{errors[status_code]}'
            }
        )


class UnknownApiError(GenericError):
    def __init__(self,
                 # headers: Mapping[str, str],
                 status_code: int = status.HTTP_500_INTERNAL_SERVER_ERROR,
                 message: str = 'Internal server error',
                 user_message: str | None = None):
        super().__init__(
            status_code=status_code,
            # headers=headers,
            message=message,
            user_message=user_message
        )


class NoContentError(Response):
    def __init__(self):
        super().__init__(
            status_code=status.HTTP_204_NO_CONTENT,
            # headers=headers,
            content=None
        )


class NotFoundError(GenericError):
    def __init__(self,
                 # headers: Mapping[str, str],
                 status_code: int = status.HTTP_404_NOT_FOUND,
                 message: str = 'Not Found',
                 user_message: str | None = None):
        super().__init__(
            status_code=status_code,
            # headers=headers,
            message=message,
            user_message=user_message
        )


class ServiceUnavailableError(GenericError):
    def __init__(self,
                 # headers: Mapping[str, str],
                 status_code: int = status.HTTP_503_SERVICE_UNAVAILABLE,
                 message: str = 'Service unavailable',
                 user_message: str | None = None):
        super().__init__(
            status_code=status_code,
            # headers=headers,
            message=message,
            user_message=user_message
        )
