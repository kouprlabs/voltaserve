# Copyright (c) 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

from typing import Dict, Iterable, Tuple

from psycopg import DatabaseError

from ..dependencies import conn
from ..errors import EmptyDataException, NotFoundException
from . import exists


def fetch_user_organizations(user_id: str, page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            if not exists(curs=curs, tablename="user", _id=user_id):
                raise NotFoundException(message=f"User with id={user_id} does not exist!")
            data = curs.execute(
                f"""
                SELECT u.id, u."permission", u.create_time as "createTime", o.id as "organizationId", 
                o."name" as "organizationName" 
                FROM userpermission u 
                JOIN organization o ON u.resource_id = o.id 
                WHERE u.user_id = '{user_id}' 
                ORDER BY u.create_time 
                OFFSET {(page - 1) * size} 
                LIMIT {size}
                """
            ).fetchall()
            if data is None or data == {}:
                raise EmptyDataException
            count = curs.execute(
                f"""
                SELECT count(1) 
                FROM userpermission u 
                JOIN organization o ON u.resource_id = o.id 
                WHERE u.user_id = '{user_id}'
                """
            ).fetchone()
            return data, count["count"]
    except DatabaseError as error:
        raise error


def fetch_user_count() -> Dict:
    try:
        with conn.cursor() as curs:
            return curs.execute('SELECT count(id) FROM "user"').fetchone()
    except DatabaseError as error:
        raise error


def fetch_user_workspaces(user_id: str, page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            if not exists(curs=curs, tablename="user", _id=user_id):
                raise NotFoundException(message=f"User with id={user_id} does not exist!")
            data = curs.execute(
                f"""
                SELECT u.id, u.permission, u.create_time as "createTime", w.id as "workspaceId", 
                w."name" as "workspaceName" 
                FROM userpermission u 
                JOIN workspace w ON u.resource_id = w.id 
                WHERE u.user_id = '{user_id}' 
                ORDER BY u.create_time 
                OFFSET {(page - 1) * size} 
                LIMIT {size}
                """
            ).fetchall()
            if data is None or data == {}:
                raise EmptyDataException
            count = curs.execute(
                f"""
                SELECT count(1) 
                FROM userpermission u 
                JOIN workspace w ON u.resource_id = w.id 
                WHERE u.user_id = '{user_id}'
                """
            ).fetchone()

            return data, count["count"]
    except DatabaseError as error:
        raise error


def fetch_user_groups(user_id: str, page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            if not exists(curs=curs, tablename="user", _id=user_id):
                raise NotFoundException(message=f"User with id={user_id} does not exist!")
            data = curs.execute(
                f"""
                SELECT u.id, u.permission, u.create_time as "createTime", g.id as "groupId", g."name" as "groupName" 
                FROM userpermission u JOIN "group" g ON u.resource_id = g.id 
                WHERE u.user_id = '{user_id}' 
                ORDER BY u.create_time 
                OFFSET {(page - 1) * size} 
                LIMIT {size}
                """
            ).fetchall()
            if data is None or data == {}:
                raise EmptyDataException
            count = curs.execute(
                f"""
                SELECT count(1) 
                FROM userpermission u 
                JOIN "group" g ON u.resource_id = g.id 
                WHERE u.user_id = '{user_id}'
                """
            ).fetchone()
            return data, count["count"]
    except DatabaseError as error:
        raise error
