# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

from psycopg2 import extras, connect, DatabaseError

from ..dependencies import settings

conn = connect(host=settings.db_host,
               user=settings.db_user,
               password=settings.db_password,
               dbname=settings.db_name,
               port=settings.db_port)


# --- FETCH --- #
def fetch_task(_id: str):
    with conn.cursor(cursor_factory=extras.RealDictCursor) as curs:
        try:
            curs.execute(f"SELECT id, name, error, percentage, is_complete, is_indeterminate, user_id, status, "
                         f"payload, task_id, create_time, update_time "
                         f"FROM {settings.db_name}.task "
                         f"WHERE id='{_id}'")
            return curs.fetchone()
        except (Exception, DatabaseError) as error:
            print(error)


def fetch_tasks(page=1, size=10):
    with conn.cursor(cursor_factory=extras.RealDictCursor) as curs:
        try:
            curs.execute(f"SELECT id, name, error, percentage, is_complete, is_indeterminate, user_id, status, "
                         f"payload, task_id, create_time, update_time "
                         f"FROM {settings.db_name}.task "
                         f"ORDER BY create_time "
                         f"OFFSET {(page - 1) * size} "
                         f"LIMIT {size}")
            data = curs.fetchall()

            curs.execute(f"SELECT count(1) "
                         f"FROM {settings.db_name}.task")
            count = curs.fetchone()

            return data, count
        except (Exception, DatabaseError) as error:
            print(error)

# --- UPDATE --- #

# --- CREATE --- #

# --- DELETE --- #
