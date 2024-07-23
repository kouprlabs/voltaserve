from typing import Annotated

from fastapi import APIRouter, Depends, status

from admin.api.database.snapshot import fetch_snapshot, fetch_snapshots
from admin.api.models.generic import GenericNotFoundResponse
from admin.api.models.snapshot import SnapshotResponse, SnapshotListRequest, SnapshotListResponse, SnapshotRequest

snapshot_api_router = APIRouter(
    prefix='/snapshot',
    tags=['snapshot'],
    responses={
        status.HTTP_404_NOT_FOUND: {
            'model': GenericNotFoundResponse
        }
    }
)


# --- GET --- #
@snapshot_api_router.get(path="/",
                         responses={
                             status.HTTP_200_OK: {
                                 'model': SnapshotResponse
                             }}
                         )
async def get_snapshot(data: Annotated[SnapshotRequest, Depends()]):
    snapshot = fetch_snapshot(_id=data.id)
    if snapshot is None:
        return GenericNotFoundResponse(message=f'Snapshot with id={data.id} does not exist')

    return SnapshotResponse(**snapshot)


@snapshot_api_router.get(path="/all",
                         responses={
                             status.HTTP_200_OK: {
                                 'model': SnapshotListResponse
                             }
                         }
                         )
async def get_all_snapshots(data: Annotated[SnapshotListRequest, Depends()]):
    snapshots = fetch_snapshots(page=data.page, size=data.size)
    if snapshots is None:
        return GenericNotFoundResponse(message=f'This instance has no snapshots')

    return SnapshotListResponse(snapshots=snapshots)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
