import datetime
from typing import List

from .generic import GenericPaginationRequest, GenericResponse, GenericListResponse, GenericRequest


# --- REQUEST MODELS --- #
class TaskRequest(GenericRequest):
    pass


class TaskListRequest(GenericPaginationRequest):
    pass


# --- RESPONSE MODELS --- #
class TaskResponse(GenericResponse):
    name: str
    error: str
    percentage: float
    is_complete: bool
    is_indeterminate: bool
    user_id: str
    status: str
    payload: str
    task_id: str
    create_time: datetime.datetime
    update_time: datetime.datetime


class TaskListResponse(GenericListResponse):
    tasks: List[TaskResponse]
