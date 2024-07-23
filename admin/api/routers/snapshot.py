from typing import Annotated

from fastapi import APIRouter, Depends, status

from admin.api.database.snapshot import fetch_snapshot, fetch_snapshots
from admin.api.models.snapshot import SnapshotResponse, SnapshotListRequest, SnapshotListResponse, SnapshotRequest

snapshot_api_router = APIRouter(
    prefix='/snapshot',
    tags=['snapshot']
)


# --- GET --- #
@snapshot_api_router.get(path="/",
                         responses={
                             status.HTTP_200_OK: {
                                 'model': SnapshotResponse
                             }}
                         )
async def get_user(data: Annotated[SnapshotRequest, Depends()]):
    return fetch_snapshot(_id=data.id)


@snapshot_api_router.get(path="/all",
                         responses={
                             status.HTTP_200_OK: {
                                 'model': SnapshotListResponse
                             }
                         }
                         )
async def get_all_users(data: Annotated[SnapshotListRequest, Depends()]):
    return fetch_snapshots(page=data.page, size=data.size)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
