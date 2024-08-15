# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.
from psycopg import DatabaseError

from ..dependencies import conn


# --- FETCH --- #
def fetch_snapshot(_id: str):
    try:
        with conn.cursor() as curs:
            curs.execute(
                f"SELECT id, version, original, preview, text, ocr, entities, mosaic, thumbnail, language, "
                f"status, task_id, create_time, update_time "
                f"FROM snapshot "
                f"WHERE id='{_id}'")
            return curs.fetchone()
    except DatabaseError as error:
        raise error


def fetch_snapshots(page=1, size=10):
    try:
        with conn.cursor() as curs:
            curs.execute(
                f"SELECT id, version, original, preview, text, ocr, entities, mosaic, thumbnail, language, "
                f"status, task_id, create_time, update_time "
                f"FROM snapshot "
                f"ORDER BY create_time "
                f"OFFSET {(page - 1) * size} "
                f"LIMIT {size}")
            data = curs.fetchall()

            curs.execute(f"SELECT count(1) "
                         f"FROM snapshot")
            count = curs.fetchone()

            return data, count
    except DatabaseError as error:
        raise error

# --- UPDATE --- #

# --- CREATE --- #

# --- DELETE --- #
