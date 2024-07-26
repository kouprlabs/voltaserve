from psycopg2 import extras, connect, DatabaseError

from ..dependencies import settings

conn = connect(host=settings.db_host,
               user=settings.db_user,
               password=settings.db_password,
               dbname=settings.db_name,
               port=settings.db_port)


# --- FETCH --- #
def fetch_user(_id: str):
    with conn.cursor(cursor_factory=extras.RealDictCursor) as curs:
        try:
            curs.execute(f"SELECT id, full_name, username, email, picture, "
                         f"is_email_confirmed, create_time, update_time "
                         f"FROM {settings.db_name}.user "
                         f"WHERE id='{_id}'")
            return curs.fetchone()
        except (Exception, DatabaseError) as error:
            print(error)


def fetch_users(page=1, size=10):
    with conn.cursor(cursor_factory=extras.RealDictCursor) as curs:
        try:
            curs.execute(f"SELECT id, full_name, username, email, picture, "
                         f"is_email_confirmed, create_time, update_time "
                         f"FROM {settings.db_name}.user "
                         f"ORDER BY create_time "
                         f"OFFSET {(page - 1) * size} "
                         f"LIMIT {(page - 1) * size + size}")
            data = curs.fetchall()

            curs.execute(f"SELECT count(1) "
                         f"FROM {settings.db_name}.user")
            count = curs.fetchone()

            return data, count
        except (Exception, DatabaseError) as error:
            print(error)


def fetch_user_organizations(user_id: str, page=1, size=10):
    with conn.cursor(cursor_factory=extras.RealDictCursor) as curs:
        try:
            curs.execute(f"SELECT * from {settings.db_name}.userpermission u "
                         f"JOIN {settings.db_name}.organization o ON u.resource_id = o.id "
                         f"WHERE u.user_id = '{user_id}' "
                         f"ORDER BY o.create_time "
                         f"OFFSET {(page - 1) * size} "
                         f"LIMIT {(page - 1) * size + size}")
            data = curs.fetchall()

            curs.execute(f"SELECT count(1) "
                         f"FROM {settings.db_name}.user "
                         f"WHERE u.user_id = '{user_id}' ")
            count = curs.fetchone()

            return data, count
        except (Exception, DatabaseError) as error:
            print(error)

# --- UPDATE --- #

# --- CREATE --- #

# --- DELETE --- #
