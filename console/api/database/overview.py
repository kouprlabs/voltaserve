# Copyright (c) 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

from psycopg import DatabaseError

from ..dependencies import conn


def fetch_version() -> str:
    try:
        with conn.cursor() as curs:
            return curs.execute("SELECT version();").fetchone()["version"].split(" ")[1]
    except DatabaseError as error:
        raise error
