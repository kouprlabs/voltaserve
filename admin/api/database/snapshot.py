from psycopg2 import extras, connect, DatabaseError

from ..dependencies import settings

conn = connect(host=settings.db_host,
               user=settings.db_user,
               password=settings.db_password,
               dbname=settings.db_name,
               port=settings.db_port)


# --- FETCH --- #
def fetch_snapshot(_id: str):
    with conn.cursor(cursor_factory=extras.RealDictCursor) as curs:
        try:
            curs.execute(
                f"SELECT id, version, original, preview, text, ocr, entities, mosaic, thumbnail, language, "
                f"status, task_id, create_time, update_time "
                f"FROM {settings.db_name}.snapshot "
                f"WHERE id='{_id}'")
            return curs.fetchone()
        except (Exception, DatabaseError) as error:
            print(error)


def fetch_snapshots(page=1, size=10):
    with conn.cursor(cursor_factory=extras.RealDictCursor) as curs:
        try:
            curs.execute(
                f"SELECT id, version, original, preview, text, ocr, entities, mosaic, thumbnail, language, "
                f"status, task_id, create_time, update_time "
                f"FROM {settings.db_name}.snapshot "
                f"ORDER BY create_time "
                f"OFFSET {(page - 1) * size} "
                f"LIMIT {(page - 1) * size + size}")
            return curs.fetchall()
        except (Exception, DatabaseError) as error:
            print(error)

# --- UPDATE --- #

# --- CREATE --- #

# --- DELETE --- #
