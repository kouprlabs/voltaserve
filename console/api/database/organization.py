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


def fetch_organization(organization_id: str) -> Dict:
    try:
        with conn.cursor() as curs:
            if not exists(curs=curs, tablename="organization", _id=organization_id):
                raise NotFoundException(message=f"Organization with id={organization_id} does not exist!")
            return curs.execute(
                f"""
                SELECT id, name, create_time as "createTime", update_time as "updateTime" 
                FROM organization 
                WHERE id='{organization_id}'
                """
            ).fetchone()
    except DatabaseError as error:
        raise error


def fetch_organizations(user_id: str, page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            data = curs.execute(
                f"""
                SELECT o.id, o.name, o.create_time, o.update_time, up.permission 
                FROM "organization" o
                LEFT JOIN  userpermission up ON up.resource_id = o.id AND up.user_id = '{user_id}' 
                ORDER BY o.create_time 
                OFFSET {(page - 1) * size} 
                LIMIT {size}
                """
            ).fetchall()
            if data is None or data == {}:
                raise EmptyDataException
            count = curs.execute('SELECT count(1) FROM "organization"').fetchone()
            return (
                {
                    "id": d.get("id"),
                    "name": d.get("name"),
                    "permission": d.get("permission"),
                    "createTime": d.get("create_time"),
                    "updateTime": d.get("update_time"),
                }
                for d in data
            ), count["count"]
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
            data = curs.execute(
                f"""
                SELECT u.user_id as "id", u2.username, u2.picture, u."permission", u2.create_time as "createTime" 
                FROM organization o 
                JOIN userpermission u ON o.id = u.resource_id 
                JOIN "user" u2 ON u.user_id = u2.id 
                WHERE o.id = '{organization_id}' 
                ORDER BY o.create_time 
                OFFSET {(page - 1) * size} 
                LIMIT {size}
                """
            ).fetchall()
            if data is None or data == {}:
                raise EmptyDataException
            count = curs.execute(
                f"""
                SELECT count(u.id) 
                FROM userpermission u 
                JOIN organization o 
                ON u.resource_id = o.id 
                WHERE o.id = '{organization_id}'
                """
            ).fetchone()
            return data, count["count"]
    except DatabaseError as error:
        raise error


def fetch_organization_workspaces(organization_id: str, page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            data = curs.execute(
                f"""
                SELECT id, name, create_time as "createTime" 
                FROM "workspace" 
                WHERE organization_id = '{organization_id}' 
                ORDER BY create_time 
                OFFSET {(page - 1) * size} 
                LIMIT {size}
                """
            ).fetchall()
            if data is None or data == {}:
                raise EmptyDataException
            count = curs.execute(
                f"""
                SELECT count(u.id) 
                FROM "userpermission" u 
                JOIN "workspace" w ON u.resource_id = w.id 
                WHERE w.organization_id = '{organization_id}'
                """
            ).fetchone()
            return data, count["count"]
    except DatabaseError as error:
        raise error


def fetch_organization_groups(organization_id: str, page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            data = curs.execute(
                f"""
                SELECT id, name, create_time as "createTime" 
                FROM "group" 
                WHERE organization_id = '{organization_id}' 
                ORDER BY create_time 
                OFFSET {(page - 1) * size} 
                LIMIT {size}
                """
            ).fetchall()
            if data is None or data == {}:
                raise EmptyDataException
            count = curs.execute(
                f"""
                SELECT count(id) 
                FROM "group" 
                WHERE organization_id = '{organization_id}'
                """
            ).fetchone()
            return data, count["count"]
    except DatabaseError as error:
        raise error
