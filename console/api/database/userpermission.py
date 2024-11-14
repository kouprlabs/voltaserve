# Copyright 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

from psycopg import DatabaseError

from api.dependencies import conn, new_id, new_timestamp


def grant_user_permission(user_id: str, resource_id: str, permission: str) -> None:
    try:
        with conn.cursor() as curs:
            curs.execute(
                """
                INSERT INTO userpermission (id, user_id, resource_id, permission, create_time)
                VALUES (%s, %s, %s, %s, %s) ON CONFLICT (user_id, resource_id) DO UPDATE SET permission = %s
                """,
                (
                    new_id(),
                    user_id,
                    resource_id,
                    permission,
                    new_timestamp(),
                    permission,
                ),
            )
    except DatabaseError as error:
        raise error


def revoke_user_permission(user_id: str, resource_id: str) -> None:
    try:
        with conn.cursor() as curs:
            curs.execute(
                "DELETE FROM userpermission WHERE user_id = %s AND resource_id = %s",
                (user_id, resource_id),
            )
    except DatabaseError as error:
        raise error
