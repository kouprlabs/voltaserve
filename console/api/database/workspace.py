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


def fetch_workspace(_id: str, user_id: str) -> Dict:
    try:
        with conn.cursor() as curs:
            if not exists(curs=curs, tablename="workspace", _id=_id):
                raise NotFoundException(message=f"Workspace with id={_id} does not exist!")
            data = curs.execute(
                f"""
                SELECT w.id, w.name, w.organization_id, o.name as "organizationName", 
                o.create_time as "organization_create_time", o.update_time as "organization_update_time", 
                w.storage_capacity, w.root_id, w.bucket, w.create_time, w.update_time, 
                up.permission 
                FROM workspace w 
                JOIN organization o ON w.organization_id = o.id 
                LEFT JOIN  userpermission up ON up.resource_id = w.id AND up.user_id = '{user_id}'
                WHERE w.id = '{_id}'
                """
            ).fetchone()
            return {
                "id": data.get("id"),
                "createTime": data.get("create_time"),
                "updateTime": data.get("update_time"),
                "name": data.get("name"),
                "storageCapacity": data.get("storage_capacity"),
                "rootId": data.get("root_id"),
                "bucket": data.get("bucket"),
                "organization": {
                    "id": data.get("organization_id"),
                    "name": data.get("organizationName"),
                    "createTime": data.get("organization_create_time"),
                    "updateTime": data.get("organization_update_time"),
                },
                "permission": data.get("permission"),
            }
    except DatabaseError as error:
        raise error


def fetch_workspace_count() -> Dict:
    try:
        with conn.cursor() as curs:
            return curs.execute('SELECT count(id) FROM "workspace"').fetchone()
    except DatabaseError as error:
        raise error


def fetch_workspaces(user_id: str, page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            data = curs.execute(
                f"""
                SELECT w.id, w.name, w.organization_id, o.name as "organization_name", 
                o.create_time as "organization_create_time", o.update_time as "organization_update_time", 
                w.storage_capacity, w.root_id, w.bucket, w.create_time, w.update_time, 
                up.permission 
                FROM workspace w 
                JOIN organization o ON w.organization_id = o.id 
                LEFT JOIN  userpermission up ON up.resource_id = w.id AND up.user_id = '{user_id}' 
                ORDER BY w.create_time 
                OFFSET {(page - 1) * size} 
                LIMIT {size}
                """
            ).fetchall()
            if data is None or data == {}:
                raise EmptyDataException
            count = curs.execute("SELECT count(1) FROM workspace").fetchone()
            return (
                {
                    "id": d.get("id"),
                    "createTime": d.get("create_time"),
                    "updateTime": d.get("update_time"),
                    "name": d.get("name"),
                    "storageCapacity": d.get("storage_capacity"),
                    "rootId": d.get("root_id"),
                    "bucket": d.get("bucket"),
                    "organization": {
                        "id": d.get("organization_id"),
                        "name": d.get("organization_name"),
                        "createTime": d.get("organization_create_time"),
                        "updateTime": d.get("organization_update_time"),
                    },
                    "permission": d.get("permission"),
                }
                for d in data
            ), count["count"]
    except DatabaseError as error:
        raise error
