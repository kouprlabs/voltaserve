import psycopg2
from psycopg2.extras import RealDictCursor

from ..dependencies import settings

conn = psycopg2.connect(host=settings.db_host,
                        user=settings.db_user,
                        password=settings.db_password,
                        dbname=settings.db_name,
                        port=settings.db_port)


# --- FETCH --- #
def fetch_task(_id: str):
    with conn.cursor(cursor_factory=psycopg2.extras.RealDictCursor) as curs:
        try:
            curs.execute(f"SELECT id, name, error, percentage, is_complete, is_indeterminate, user_id, status, payload, task_id, create_time, update_time "
                         f"FROM {settings.db_name}.task "
                         f"WHERE id='{_id}'")
            return curs.fetchone()
        except (Exception, psycopg2.DatabaseError) as error:
            print(error)


def fetch_tasks(page=0, size=10):
    with conn.cursor(cursor_factory=psycopg2.extras.RealDictCursor) as curs:
        try:
            curs.execute(f"SELECT id, name, error, percentage, is_complete, is_indeterminate, user_id, status, payload, task_id, create_time, update_time "
                         f"FROM {settings.db_name}.task "
                         f"ORDER BY create_time "
                         f"OFFSET {page * size} "
                         f"LIMIT {page * size + size}")
            return curs.fetchall()
        except (Exception, psycopg2.DatabaseError) as error:
            print(error)

# --- UPDATE --- #

# --- CREATE --- #

# --- DELETE --- #

