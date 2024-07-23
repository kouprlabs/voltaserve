import datetime
from typing import List

from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class SnapshotFileRequest(GenericRequest):
    pass


class SnapshotFileListRequest(GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class SnapshotFileResponse(GenericResponse):
    snapshot_id: str
    file_id: str
    create_time: datetime.datetime


class SnapshotFileListResponse(GenericListResponse):
    snapshotfiles: List[SnapshotFileResponse]
