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


def fetch_invitation(_id: str) -> Dict | None:
    try:
        with conn.cursor() as curs:
            if not exists(curs=curs, tablename="invitation", _id=_id):
                raise NotFoundException(
                    message=f"Invitation with id={_id} does not exist!"
                )

            data = curs.execute(
                f"""
                SELECT id, organization_id, owner_id, email, status, create_time as "createTime", 
                update_time as "updateTime", o.name as "organization_name", 
                o.create_time as "organization_create_time", o.update_time as "organization_update_time", 
                i.owner_id as "ownerId" 
                FROM invitation i 
                JOIN organization o ON i.organization_id = o.id 
                WHERE i.id=\'{_id}\'
                """
            ).fetchone()
            return data if data != {} else None
    except DatabaseError as error:
        raise error


def fetch_invitations(page=1, size=10) -> Tuple[Iterable[Dict], int]:
    try:
        with conn.cursor() as curs:
            data = curs.execute(
                f"""
                SELECT i.id, i.organization_id, o.name as "organization_name", 
                o.create_time as "organization_create_time", o.update_time as "organization_update_time", 
                i.owner_id, i.email, i.status, i.create_time, i.update_time 
                FROM invitation i 
                JOIN organization o ON i.organization_id = o.id 
                ORDER BY i.create_time 
                OFFSET {(page - 1) * size} 
                LIMIT {size}
                """
            ).fetchall()

            if data is None or data == {}:
                raise EmptyDataException

            count = curs.execute("SELECT count(1) FROM invitation").fetchone()

            return (
                {
                    "id": d.get("id"),
                    "ownerId": d.get("owner_id"),
                    "createTime": d.get("create_time"),
                    "updateTime": d.get("update_time"),
                    "email": d.get("email"),
                    "status": d.get("status"),
                    "organization": {
                        "id": d.get("organization_id"),
                        "name": d.get("organization_name"),
                        "createTime": d.get("organization_create_time"),
                        "updateTime": d.get("organization_update_time"),
                    },
                }
                for d in data
            ), count["count"]
    except DatabaseError as error:
        raise error


def update_invitation(data: dict) -> None:
    try:
        data = {
            "id": data.get("id"),
            "status": "accepted" if data.get("accept") else "declined",
        }
        with conn.cursor() as curs:
            curs.execute(parse_sql_update_query("invitation", data))
    except DatabaseError as error:
        raise error
