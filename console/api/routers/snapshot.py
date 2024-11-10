# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

from typing import Annotated

from fastapi import APIRouter, Depends, status

from ..database import fetch_snapshot, fetch_snapshots
from ..dependencies import JWTBearer
from ..log import base_logger
from ..errors import (
    NotFoundError,
    EmptyDataException,
    NoContentError,
    NotFoundException,
    UnknownApiError,
)
from ..models import (
    SnapshotResponse,
    SnapshotListRequest,
    SnapshotListResponse,
    SnapshotRequest,
)

snapshot_api_router = APIRouter(
    prefix="/snapshot", tags=["snapshot"], dependencies=[Depends(JWTBearer())]
)


logger = base_logger.getChild("snapshot")


# --- GET --- #
@snapshot_api_router.get(
    path="", responses={status.HTTP_200_OK: {"model": SnapshotResponse}}
)
async def get_snapshot(data: Annotated[SnapshotRequest, Depends()]):
    try:
        snapshot = fetch_snapshot(_id=data.id)
        if snapshot is None:
            return NotFoundError(message=f"Snapshot with id={data.id} does not exist")

        return SnapshotResponse(**snapshot)
    except NotFoundException as e:
        logger.error(e)
        return NotFoundError(message=str(e))
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


@snapshot_api_router.get(
    path="/all", responses={status.HTTP_200_OK: {"model": SnapshotListResponse}}
)
async def get_all_snapshots(data: Annotated[SnapshotListRequest, Depends()]):
    try:
        snapshots, count = fetch_snapshots(page=data.page, size=data.size)

        return SnapshotListResponse(
            data=snapshots,
            totalElements=count,
            totalPages=(count + data.size - 1) // data.size,
            page=data.page,
            size=data.size,
        )
    except EmptyDataException as e:
        logger.error(e)
        return NoContentError()
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()


# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
