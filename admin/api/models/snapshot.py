import datetime
from typing import List

from pydantic import EmailStr

from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class SnapshotRequest(GenericRequest):
    pass


class SnapshotListRequest(GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class SnapshotResponse(GenericResponse):
    version: str
    original: dict
    preview: dict
    text: dict
    ocr: dict
    entities: dict
    mosaic: dict
    thumbnail: dict
    language: str
    status: str
    task_id: str
    create_time: datetime.datetime
    update_time: datetime.datetime


class SnapshotListResponse(GenericListResponse):
    snapshots: List[SnapshotResponse]
