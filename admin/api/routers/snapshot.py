from typing import Annotated

from fastapi import APIRouter, Depends, status

from ..database import fetch_snapshot, fetch_snapshots
from ..exceptions import GenericNotFoundException
from ..models import GenericNotFoundResponse, SnapshotResponse, SnapshotListRequest, SnapshotListResponse, \
    SnapshotRequest

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
        raise GenericNotFoundException(detail=f'Snapshot with id={data.id} does not exist')

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
        raise GenericNotFoundException(detail='This instance has no snapshots')

    return SnapshotListResponse(snapshots=snapshots)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
