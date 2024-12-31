# Copyright (c) 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

from psycopg import connect
from psycopg.rows import dict_row

from . import settings

conn = connect(
    conninfo=f"postgres://{settings.POSTGRES_USER}:"
    f"{settings.POSTGRES_PASSWORD}@"
    f"{settings.POSTGRES_URL}:"
    f"{settings.POSTGRES_PORT}/"
    f"{settings.POSTGRES_NAME}",
    row_factory=dict_row,
    autocommit=True,
)
