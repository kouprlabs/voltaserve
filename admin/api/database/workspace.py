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
def fetch_workspace(_id: str) -> Dict:
    try:
        with conn.cursor() as curs:
            if not exists(curs=curs, tablename='workspace', _id=_id):
                raise NotFoundException(message=f'Workspace with id={_id} does not exist!')

            curs.execute(
                f'SELECT w.id, w.name, w.organization_id as "organizationId", o.name as "organizationName", '
                f'w.storage_capacity as "storageCapacity", w.root_id as "rootId", w.bucket,'
                f' w.create_time as "createTime", w.update_time as "updateTime" '
                f'FROM workspace w join organization o on w.organization_id = o.id '
                f'WHERE w.id = \'{_id}\'')
            return curs.fetchone()
    except DatabaseError as error:
        raise error


def fetch_workspaces(page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            data = curs.execute(
                f'SELECT w.id, w.name, w.organization_id, o.name as "organization_name", '
                f'o.create_time as "organization_create_time", o.update_time as "organization_update_time", '
                f'w.storage_capacity, w.root_id, w.bucket, w.create_time, w.update_time '
                f'FROM workspace w join organization o on w.organization_id = o.id '
                f'ORDER BY w.create_time '
                f'OFFSET {(page - 1) * size} '
                f'LIMIT {size}').fetchall()

            if data is None or data == {}:
                raise EmptyDataException

            count = curs.execute('SELECT count(1) FROM workspace').fetchone()

            return ({'id': d.get('id'),
                     'createTime': d.get('create_time'),
                     'updateTime': d.get('update_time'),
                     'name': d.get('name'),
                     'storageCapacity': d.get('storage_capacity'),
                     'rootId': d.get('root_id'),
                     'bucket': d.get('bucket'),
                     'organization': {
                         'id': d.get('organization_id'),
                         'name': d.get('organization_name'),
                         'createTime': d.get('organization_create_time'),
                         'updateTime': d.get('organization_update_time')}
                     } for d in data), count['count']

    except DatabaseError as error:
        raise error


# --- UPDATE --- #
def update_workspace(data: dict) -> None:
    try:
        with conn.cursor() as curs:
            curs.execute(parse_sql_update_query('workspace', data))
    except DatabaseError as error:
        raise error

# --- CREATE --- #

# --- DELETE --- #
