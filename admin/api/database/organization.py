# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.
from typing import Dict, Tuple, Iterable

from psycopg import DatabaseError

from . import exists
from ..dependencies import conn, parse_sql_update_query
from ..errors import EmptyDataException, NotFoundException


# --- FETCH --- #
def fetch_organization(organization_id: str) -> Dict:
    try:
        with conn.cursor() as curs:
            if not exists(curs=curs, tablename='organization', _id=organization_id):
                raise NotFoundException(message=f'Organization with id={organization_id} does not exist!')

            return curs.execute(f'SELECT id, name, create_time as "createTime", update_time as "updateTime" '
                                f'FROM organization '
                                f'WHERE id=\'{organization_id}\'').fetchone()

    except DatabaseError as error:
        raise error


def fetch_organizations(page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            data = curs.execute(f'SELECT id, name, create_time as "createTime", update_time as "updateTime" '
                                f'FROM "organization" '
                                f'ORDER BY create_time '
                                f'OFFSET {(page - 1) * size} '
                                f'LIMIT {size}').fetchall()

            if data is None or data == {}:
                raise EmptyDataException

            count = curs.execute('SELECT count(1) FROM "organization"').fetchone()

            return data, count['count']

    except DatabaseError as error:
        raise error


def fetch_organization_count() -> Dict:
    try:
        with conn.cursor() as curs:
            return curs.execute('SELECT count(id) FROM "organization"').fetchone()

    except DatabaseError as error:
        raise error


def fetch_organization_users(organization_id: str, page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            data = curs.execute(f'SELECT u.user_id as "id", u2.username, u2.picture, u."permission", '
                                f'u.create_time as "createTime" '
                                f'FROM organization o '
                                f'JOIN userpermission u on o.id = u.resource_id '
                                f'JOIN "user" u2 on u.user_id = u2.id '
                                f'WHERE o.id = \'{organization_id}\' '
                                f'ORDER BY o.create_time '
                                f'OFFSET {(page - 1) * size} '
                                f'LIMIT {size}').fetchall()

            if data is None or data == {}:
                raise EmptyDataException

            return data, len(data)

    except DatabaseError as error:
        raise error


def fetch_organization_workspaces(organization_id: str, page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            data = curs.execute(
                f'SELECT id, name, create_time as "createTime"'
                f'FROM workspace '
                f'WHERE organization_id = \'{organization_id}\' '
                f'ORDER BY create_time '
                f'OFFSET {(page - 1) * size} '
                f'LIMIT {size}').fetchall()

            if data is None or data == {}:
                raise EmptyDataException

            return data, len(data)
    except DatabaseError as error:
        raise error


def fetch_organization_groups(organization_id: str, page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            data = curs.execute(
                f'SELECT id, name, create_time as "createTime"'
                f'FROM "group" '
                f'WHERE organization_id = \'{organization_id}\' '
                f'ORDER BY create_time '
                f'OFFSET {(page - 1) * size} '
                f'LIMIT {size}').fetchall()

            if data is None or data == {}:
                raise EmptyDataException

            return data, len(data)
    except DatabaseError as error:
        raise error


# --- UPDATE --- #
def update_organization(data: Dict) -> Dict:
    try:
        with conn.cursor() as curs:
            curs.execute(parse_sql_update_query('organization', data))
    except DatabaseError as error:
        raise error

# --- CREATE --- #

# --- DELETE --- #
