# Copyright (c) 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.


class GenericDatabaseException(Exception):
    def __init__(self, message: str):
        super().__init__(message)


class NotFoundException(GenericDatabaseException):
    def __init__(self, message: str = "ID does not exist"):
        super().__init__(message)


class EmptyDataException(GenericDatabaseException):
    def __init__(self, message: str = "Query returned no data"):
        super().__init__(message)
