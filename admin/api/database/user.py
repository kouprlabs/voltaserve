# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.
from typing import Tuple, Iterable, Dict

from psycopg import DatabaseError

from . import exists
from ..dependencies import conn
from ..errors import EmptyDataException, NotFoundException


# --- FETCH --- #

def fetch_user_organizations(user_id: str, page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            if not exists(curs=curs, tablename='user', _id=user_id):
                raise NotFoundException(message=f'User with id={user_id} does not exist!')

            data = curs.execute(f'SELECT u.id, u."permission", u.create_time as "createTime", '
                                f'o.id as "organizationId", o."name" as "organizationName" from userpermission u '
                                f'JOIN organization o ON u.resource_id = o.id '
                                f'WHERE u.user_id = \'{user_id}\' '
                                f'ORDER BY u.create_time '
                                f'OFFSET {(page - 1) * size} '
                                f'LIMIT {size}').fetchall()

            if data is None or data == {}:
                raise EmptyDataException

            count = curs.execute(f'SELECT count(1) '
                                 f'FROM userpermission u '
                                 f'JOIN organization o ON u.resource_id = o.id '
                                 f'WHERE u.user_id = \'{user_id}\'').fetchone()

            return data, count['count']
    except DatabaseError as error:
        raise error


def fetch_user_workspaces(user_id: str, page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            if not exists(curs=curs, tablename='user', _id=user_id):
                raise NotFoundException(message=f'User with id={user_id} does not exist!')

            data = curs.execute(f'SELECT u.id, u.permission, u.create_time as "createTime", '
                                f'w.id as "workspaceId", w."name" as "workspaceName" '
                                f'FROM userpermission u '
                                f'JOIN workspace w ON u.resource_id = w.id WHERE u.user_id = \'{user_id}\' '
                                f'ORDER BY u.create_time '
                                f'OFFSET {(page - 1) * size} '
                                f'LIMIT {size}').fetchall()

            if data is None or data == {}:
                raise EmptyDataException

            count = curs.execute(f'SELECT count(1)'
                                 f'FROM userpermission u '
                                 f'JOIN workspace w ON u.resource_id = w.id '
                                 f'WHERE u.user_id = \'{user_id}\'').fetchone()

            return data, count['count']
    except DatabaseError as error:
        raise error


def fetch_user_groups(user_id: str, page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            if not exists(curs=curs, tablename='user', _id=user_id):
                raise NotFoundException(message=f'User with id={user_id} does not exist!')

            data = curs.execute(f'SELECT u.id, u.permission, u.create_time as "createTime", '
                                f'g.id as "groupId", g."name" as "groupName" '
                                f'FROM userpermission u '
                                f'JOIN "group" g ON u.resource_id = g.id WHERE u.user_id = \'{user_id}\' '
                                f'ORDER BY u.create_time '
                                f'OFFSET {(page - 1) * size} '
                                f'LIMIT {size}').fetchall()

            if data is None or data == {}:
                raise EmptyDataException

            count = curs.execute(f'SELECT count(1)'
                                 f'FROM userpermission u '
                                 f'JOIN "group" g ON u.resource_id = g.id '
                                 f'WHERE u.user_id = \'{user_id}\'').fetchone()

            return data, count['count']
    except DatabaseError as error:
        raise error

# --- UPDATE --- #

# --- CREATE --- #

# --- DELETE --- #
