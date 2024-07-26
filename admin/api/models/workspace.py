import datetime
from typing import List

from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class WorkspaceRequest(GenericRequest):
    pass


class WorkspaceListRequest(GenericPaginationRequest):
    pass


class OrganizationWorkspaceListRequest(WorkspaceRequest, WorkspaceListRequest):
    pass


# --- RESPONSE MODELS --- #
class WorkspaceResponse(GenericResponse):
    name: str
    organization_id: str
    storage_capacity: float
    root_id: str
    bucket: str
    create_time: datetime.datetime
    update_time: datetime.datetime


class WorkspaceListResponse(GenericListResponse):
    data: List[WorkspaceResponse]
