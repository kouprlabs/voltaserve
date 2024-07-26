from psycopg2 import extras, connect, DatabaseError

from ..dependencies import settings

conn = connect(host=settings.db_host,
               user=settings.db_user,
               password=settings.db_password,
               dbname=settings.db_name,
               port=settings.db_port)


# --- FETCH --- #
def fetch_group(_id: str):
    with conn.cursor(cursor_factory=extras.RealDictCursor) as curs:
        try:
            curs.execute(f"SELECT id, name, organization_id, create_time, update_time "
                         f"FROM {settings.db_name}.group "
                         f"WHERE id='{_id}'")
            return curs.fetchone()
        except (Exception, DatabaseError) as error:
            print(error)


def fetch_groups(page=1, size=10):
    with conn.cursor(cursor_factory=extras.RealDictCursor) as curs:
        try:
            curs.execute(f"SELECT id, name, organization_id, create_time, update_time "
                         f"FROM {settings.db_name}.group "
                         f"ORDER BY create_time "
                         f"OFFSET {(page - 1) * size} "
                         f"LIMIT {(page - 1) * size + size}")
            x = curs.fetchall()
            return x
        except (Exception, DatabaseError) as error:
            print(error)

# --- UPDATE --- #

# --- CREATE --- #

# --- DELETE --- #
