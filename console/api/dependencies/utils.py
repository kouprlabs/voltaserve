# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.
from datetime import datetime
from functools import reduce


def camel_to_snake(data: str):
    return reduce(lambda x, y: x + ('_' if y.isupper() else '') + y, data).lower()


def parse_sql_update_query(tablename: str, data: dict):
    date_fields = ", ".join(f'{camel_to_snake(k)} = \'{v.strftime("%Y-%m-%dT%H:%M:%SZ")}\'' for
                            k, v in data.items() if k != "id" and isinstance(v, datetime))
    data_fields = ", ".join(f'{camel_to_snake(k)} = \'{v}\'' for
                            k, v in data.items() if k != "id" and not isinstance(v, datetime))
    return (f'UPDATE "{tablename}" SET {data_fields} WHERE id = \'{data["id"]}\''
            f';') if date_fields == '' else (f'UPDATE "{tablename}" SET '
                                             f'{data_fields}, {date_fields} WHERE id = \'{data["id"]}\';')