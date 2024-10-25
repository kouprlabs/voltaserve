# Copyright 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

from flask import Blueprint

bp = Blueprint("version", __name__)


@bp.route("/version", methods=["GET"])
def get_version():
    return {"version": "3.0.0"}
