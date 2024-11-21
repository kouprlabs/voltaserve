# Copyright (c) 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

from typing import Tuple, Iterable, Dict

from psycopg import DatabaseError

from ..dependencies import conn, parse_sql_update_query
from ..errors import EmptyDataException, NotFoundException
from .generic import exists


def fetch_group(_id: str) -> Dict:
    try:
        with conn.cursor() as curs:
            if not exists(curs=curs, tablename="group", _id=_id):
                raise NotFoundException(message=f"Group with id={_id} does not exist!")

            data = curs.execute(
                f"""
                SELECT g.id as "group_id", g."name" as "group_name", g.create_time, g.update_time, o.id as "org_id", 
                o."name" as "org_name", o.create_time as "org_create_time", o.update_time as "org_update_time" 
                FROM "group" g 
                JOIN organization o ON g.organization_id = o.id 
                WHERE g.id='{_id}'
                """
            ).fetchone()

            return (
                {
                    "createTime": data.get("create_time"),
                    "id": data.get("group_id"),
                    "name": data.get("group_name"),
                    "organization": {
                        "id": data.get("org_id"),
                        "name": data.get("org_name"),
                        "createTime": data.get("org_create_time"),
                        "updateTime": data.get("org_update_time"),
                    },
                }
                if data is not None
                else None
            )

    except DatabaseError as error:
        raise error


def fetch_group_count() -> Dict:
    try:
        with conn.cursor() as curs:
            return curs.execute('SELECT count(id) FROM "group"').fetchone()

    except DatabaseError as error:
        raise error


def fetch_groups(page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            data = curs.execute(
                f"""
                SELECT g.id as "group_id", g."name" as "group_name", g.create_time as "group_create_time", 
                g.update_time as "group_update_time", o.id as "org_id", o."name" as "org_name", 
                o.create_time as "org_create_time", o.update_time as "org_update_time" 
                FROM "group" g 
                JOIN organization o ON g.organization_id = o.id 
                ORDER BY g.create_time 
                OFFSET {(page - 1) * size} 
                LIMIT {size}
                """
            ).fetchall()

            if data is None or data == {}:
                raise EmptyDataException

            count = curs.execute('SELECT count(1) FROM "group"').fetchone()

            return (
                {
                    "id": d.get("group_id"),
                    "createTime": d.get("group_create_time"),
                    "updateTime": d.get("group_update_time"),
                    "name": d.get("group_name"),
                    "organization": {
                        "id": d.get("org_id"),
                        "name": d.get("org_name"),
                        "createTime": d.get("org_create_time"),
                        "updateTime": d.get("org_update_time"),
                    },
                }
                for d in data
            ), count["count"]

    except DatabaseError as error:
        raise error
