# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.
from typing import Dict, Tuple, Iterable

from psycopg import DatabaseError

from . import exists
from ..dependencies import conn
from ..errors import EmptyDataException, NotFoundException


# --- FETCH --- #
def fetch_snapshot(_id: str) -> Dict:
    try:
        with conn.cursor() as curs:
            if not exists(curs=curs, tablename="snapshot", _id=_id):
                raise NotFoundException(
                    message=f"Snapshot with id={_id} does not exist!"
                )

            return curs.execute(
                f"""
                SELECT id, version, original, preview, text, ocr, entities, mosaic, thumbnail, language, status, 
                task_id, create_time, update_time 
                FROM snapshot 
                WHERE id='{_id}'
                """
            ).fetchone()
    except DatabaseError as error:
        raise error


def fetch_snapshots(page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            data = curs.execute(
                f"""
                SELECT id, version, original, preview, text, ocr, entities, mosaic, thumbnail, language, 
                status, task_id, create_time, update_time 
                FROM snapshot 
                ORDER BY create_time 
                OFFSET {(page - 1) * size} 
                LIMIT {size}
                """
            ).fetchall()

            if data is None or data == {}:
                raise EmptyDataException

            count = curs.execute("SELECT count(1) FROM snapshot").fetchone()

            return data, count["count"]
    except DatabaseError as error:
        raise error


# --- UPDATE --- #

# --- CREATE --- #

# --- DELETE --- #
