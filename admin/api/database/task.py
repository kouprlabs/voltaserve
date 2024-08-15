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
def fetch_task(_id: str):
    try:
        with conn.cursor() as curs:
            curs.execute(f"SELECT id, name, error, percentage, is_complete, is_indeterminate, user_id, status, "
                         f"payload, task_id, create_time, update_time "
                         f"FROM task "
                         f"WHERE id='{_id}'")
            return curs.fetchone()
    except DatabaseError as error:
        raise error


def fetch_tasks(page=1, size=10):
    try:
        with conn.cursor() as curs:
            curs.execute(f"SELECT id, name, error, percentage, is_complete, is_indeterminate, user_id, status, "
                         f"payload, task_id, create_time, update_time "
                         f"FROM task "
                         f"ORDER BY create_time "
                         f"OFFSET {(page - 1) * size} "
                         f"LIMIT {size}")
            data = curs.fetchall()

            curs.execute(f"SELECT count(1) "
                         f"FROM task")
            count = curs.fetchone()

            return data, count['count']
    except DatabaseError as error:
        raise error

# --- UPDATE --- #

# --- CREATE --- #

# --- DELETE --- #
