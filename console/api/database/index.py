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

from ..dependencies import conn
from ..errors import EmptyDataException


# --- FETCH --- #
def fetch_index(_id: str) -> Dict:
    try:
        with conn.cursor() as curs:
            data = curs.execute(
                f"""
                SELECT tablename, indexname, indexdef 
                FROM pg_indexes 
                WHERE schemaname = 'public' AND indexname='{_id}'
                """
            ).fetchone()
            return {
                "tableName": data.get("tablename"),
                "indexName": data.get("indexname"),
                "indexDef": data.get("indexdef"),
            }
    except DatabaseError as error:
        raise error


def fetch_indexes(page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            data = curs.execute(
                f"""
                SELECT tablename, indexname, indexdef 
                FROM pg_indexes 
                WHERE schemaname = 'public' 
                ORDER BY tablename, indexname 
                OFFSET {(page - 1) * size} 
                LIMIT {size}
                """
            ).fetchall()

            if data is None or data == {}:
                raise EmptyDataException

            curs.execute("SELECT count(1) FROM pg_indexes WHERE schemaname = 'public' ")
            count = curs.fetchone()

            return (
                {
                    "tableName": d.get("tablename"),
                    "indexName": d.get("indexname"),
                    "indexDef": d.get("indexdef"),
                }
                for d in data
            ), count["count"]
    except DatabaseError as error:
        raise error


# --- UPDATE --- #

# --- CREATE --- #

# --- DELETE --- #
