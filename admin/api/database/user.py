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

def fetch_user_organizations(user_id: str, page=1, size=10):
    try:
        with conn.cursor() as curs:
            data = curs.execute(f'SELECT u.id, u."permission", u.create_time as "createTime", '
                                f'o.id as "organizationId", o."name" as "organizationName" from userpermission u '
                                f'JOIN organization o ON u.resource_id = o.id '
                                f'WHERE u.user_id = \'{user_id}\' '
                                f'ORDER BY u.create_time '
                                f'OFFSET {(page - 1) * size} '
                                f'LIMIT {size}').fetchall()

            count = curs.execute(f'SELECT count(1) '
                                 f'FROM userpermission u '
                                 f'JOIN organization o ON u.resource_id = o.id '
                                 f'WHERE u.user_id = \'{user_id}\'').fetchone()

            return data, count['count']
    except DatabaseError as error:
        raise error


def fetch_user_workspaces(user_id: str, page=1, size=10):
    try:
        with conn.cursor() as curs:
            data = curs.execute(f'SELECT u.id, u.permission, u.create_time as "createTime", '
                                f'w.id as "workspaceId", w."name" as "workspaceName" '
                                f'FROM userpermission u '
                                f'JOIN workspace w ON u.resource_id = w.id WHERE u.user_id = \'{user_id}\' '
                                f'ORDER BY u.create_time '
                                f'OFFSET {(page - 1) * size} '
                                f'LIMIT {size}').fetchall()

            count = curs.execute(f'SELECT count(1)'
                                 f'FROM userpermission u '
                                 f'JOIN workspace w ON u.resource_id = w.id '
                                 f'WHERE u.user_id = \'{user_id}\'').fetchone()

            return data, count['count']
    except DatabaseError as error:
        raise error


def fetch_user_groups(user_id: str, page=1, size=10):
    try:
        with conn.cursor() as curs:
            data = curs.execute(f'SELECT u.id, u.permission, u.create_time as "createTime", '
                                f'g.id as "groupId", g."name" as "groupName" '
                                f'FROM userpermission u '
                                f'JOIN "group" g ON u.resource_id = g.id WHERE u.user_id = \'{user_id}\' '
                                f'ORDER BY u.create_time '
                                f'OFFSET {(page - 1) * size} '
                                f'LIMIT {size}').fetchall()

            count = curs.execute(f'SELECT count(1)'
                                 f'FROM userpermission u '
                                 f'JOIN workspace w ON u.resource_id = w.id '
                                 f'WHERE u.user_id = \'{user_id}\'').fetchone()

            return data, count['count']
    except DatabaseError as error:
        raise error

# --- UPDATE --- #

# --- CREATE --- #

# --- DELETE --- #
