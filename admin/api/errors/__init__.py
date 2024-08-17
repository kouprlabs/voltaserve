# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

from .api_errors import NotFoundError, NoContentError, ServiceUnavailableError, \
    UnknownApiError
from .database_errors import EmptyDataException, NotFoundException
from .error_codes import errors
