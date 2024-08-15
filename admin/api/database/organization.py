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

from ..dependencies import conn, parse_sql_update_query


# --- FETCH --- #
def fetch_organization(organization_id: str):
    try:
        with conn.cursor() as curs:
            return curs.execute(f'SELECT id, name, create_time as "createTime", update_time as "updateTime" '
                                f'FROM organization '
                                f'WHERE id=\'{organization_id}\'').fetchone()

    except DatabaseError as error:
        raise error


def fetch_organizations(page=1, size=10):
    try:
        with conn.cursor() as curs:
            data = curs.execute(f'SELECT id, name, create_time as "createTime", update_time as "updateTime" '
                                f'FROM "organization" '
                                f'ORDER BY create_time '
                                f'OFFSET {(page - 1) * size} '
                                f'LIMIT {size}').fetchall()

            count = curs.execute('SELECT count(1) FROM "organization"').fetchone()

            return data, count['count']

    except DatabaseError as error:
        raise error


def fetch_organization_users(organization_id: str, page=1, size=10):
    try:
        with conn.cursor() as curs:
            data = curs.execute(f'SELECT u.user_id as "id", u2.username, u."permission", u.create_time as "createTime" '
                                f'FROM organization o '
                                f'JOIN userpermission u on o.id = u.resource_id '
                                f'JOIN "user" u2 on u.user_id = u2.id '
                                f'WHERE o.id = \'{organization_id}\' '
                                f'ORDER BY o.create_time '
                                f'OFFSET {(page - 1) * size} '
                                f'LIMIT {size}').fetchall()

            return data, len(data)

    except DatabaseError as error:
        raise error


def fetch_organization_workspaces(organization_id: str, page=1, size=10):
    try:
        with conn.cursor() as curs:
            data = curs.execute(
                f'SELECT id, name, storage_capacity as "storageCapacity", '
                f'create_time as "createTime", update_time as "updateTime" '
                f'FROM workspace '
                f'WHERE organization_id = \'{organization_id}\' '
                f'ORDER BY create_time '
                f'OFFSET {(page - 1) * size} '
                f'LIMIT {size}').fetchall()

            return data, len(data)
    except DatabaseError as error:
        raise error


# --- UPDATE --- #
def update_organization(data: dict):
    try:
        with conn.cursor() as curs:
            curs.execute(parse_sql_update_query('organization', data))
    except DatabaseError as error:
        raise error

# --- CREATE --- #

# --- DELETE --- #
