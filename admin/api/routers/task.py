from typing import Annotated

from fastapi import APIRouter, Depends, status

from admin.api.database.task import fetch_task, fetch_tasks
from admin.api.models.task import TaskResponse, TaskRequest, TaskListResponse, TaskListRequest

task_api_router = APIRouter(
    prefix='/task'
)


# --- GET --- #
@task_api_router.get(path="/",
                     responses={
                         status.HTTP_200_OK: {
                             'model': TaskResponse
                         }}
                     )
async def get_task(data: Annotated[TaskRequest, Depends()]):
    return fetch_task(_id=data.id)


@task_api_router.get(path="/all",
                     responses={
                         status.HTTP_200_OK: {
                             'model': TaskListResponse
                         }
                     }
                     )
async def get_all_tasks(data: Annotated[TaskListRequest, Depends()]):
    return fetch_tasks(page=data.page, size=data.size)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
