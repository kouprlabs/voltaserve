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
def fetch_group(_id: str):
    try:
        with conn.cursor() as curs:
            data = curs.execute(f'SELECT g.id as "group_id", g."name" as "group_name", g.create_time, g.update_time'
                                f'o.id as "org_id", o."name" as "org_name", o.create_time as "org_create_time", '
                                f'o.update_time as "org_update_time"'
                                f'FROM "group" g '
                                f'JOIN organization o ON g.organization_id = o.id '
                                f'WHERE g.id=\'{_id}\'').fetchone()

            return {'createTime': data.get('create_time'),
                    'id': data.get('group_id'),
                    'name': data.get('group_name'),
                    'organization': {
                        'id': data.get('org_id'),
                        'name': data.get('org_name'),
                        'createTime': data.get('org_create_time'),
                        'updateTime': data.get('org_update_time')
                    }
                    }
    except DatabaseError as error:
        raise error


def fetch_groups(page=1, size=10):
    try:
        with conn.cursor() as curs:
            data = curs.execute(f'SELECT g.id as "group_id", g."name" as "group_name", g.create_time '
                                f'as "group_create_time", g.update_time as "group_update_time", '
                                f'o.id as "org_id", o."name" as "org_name", o.create_time as "org_create_time", '
                                f'o.update_time as "org_update_time"'
                                f'FROM "group" g '
                                f'JOIN organization o ON g.organization_id = o.id '
                                f'ORDER BY g.create_time '
                                f'OFFSET {(page - 1) * size} '
                                f'LIMIT {size}').fetchall()

            count = curs.execute(f'SELECT count(1) FROM "group"').fetchone()

            return ({'id': d.get('group_id'),
                     'createTime': d.get('group_create_time'),
                     'updateTime': d.get('group_update_time'),
                     'name': d.get('group_name'),
                     'organization': {
                         'id': d.get('org_id'),
                         'name': d.get('org_name'),
                         'createTime': d.get('org_create_time'),
                         'updateTime': d.get('org_update_time'),
                     }
                     } for d in data), count['count']

    except DatabaseError as error:
        raise error


# --- UPDATE --- #
def update_group(data: dict):
    try:
        with conn.cursor() as curs:
            curs.execute(parse_sql_update_query('group', data))
    except DatabaseError as error:
        raise error

# --- CREATE --- #

# --- DELETE --- #
