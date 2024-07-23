import psycopg2
from psycopg2.extras import RealDictCursor

from ..dependencies import settings

conn = psycopg2.connect(host=settings.db_host,
                        user=settings.db_user,
                        password=settings.db_password,
                        dbname=settings.db_name,
                        port=settings.db_port)


# --- FETCH --- #
def fetch_user(_id: str):
    with conn.cursor(cursor_factory=psycopg2.extras.RealDictCursor) as curs:
        try:
            curs.execute(f"SELECT id, full_name, username, email, is_email_confirmed, create_time, update_time "
                         f"FROM {settings.db_name}.user")
            return curs.fetchone()
        except (Exception, psycopg2.DatabaseError) as error:
            print(error)


def fetch_users(page=0, size=10):
    with conn.cursor(cursor_factory=psycopg2.extras.RealDictCursor) as curs:
        try:
            curs.execute(f"SELECT id, full_name, username, email, is_email_confirmed, create_time, update_time "
                         f"FROM {settings.db_name}.user "
                         f"ORDER BY create_time "
                         f"OFFSET {page * size} "
                         f"LIMIT {page * size + size}")
            return curs.fetchall()
        except (Exception, psycopg2.DatabaseError) as error:
            print(error)

# --- UPDATE --- #

# --- CREATE --- #

# --- DELETE --- #

