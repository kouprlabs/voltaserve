# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

from fastapi import HTTPException, status


class GenericUnauthorizedException(HTTPException):
    def __init__(self, detail='Unauthorized'):
        super().__init__(status_code=status.HTTP_401_UNAUTHORIZED, detail=detail)


class GenericForbiddenException(HTTPException):
    def __init__(self, detail='Forbidden'):
        super().__init__(status_code=status.HTTP_403_FORBIDDEN, detail=detail)


class GenericUnexpectedException(HTTPException):
    def __init__(self, detail='Unexpected'):
        super().__init__(status_code=status.HTTP_500_INTERNAL_SERVER_ERROR, detail=detail)
