from typing import Annotated

from fastapi import APIRouter, Depends, status

from ..database import fetch_task, fetch_tasks
from ..exceptions import GenericNotFoundException
from ..models import GenericNotFoundResponse, TaskResponse, TaskRequest, TaskListResponse, TaskListRequest

task_api_router = APIRouter(
    prefix='/task',
    tags=['task'],
    responses={
        status.HTTP_404_NOT_FOUND: {
            'model': GenericNotFoundResponse
        }
    }
)


# --- GET --- #
@task_api_router.get(path="/",
                     responses={
                         status.HTTP_200_OK: {
                             'model': TaskResponse
                         }}
                     )
async def get_task(data: Annotated[TaskRequest, Depends()]):
    task = fetch_task(_id=data.id)
    if task is None:
        raise GenericNotFoundException(detail=f'Task with id={data.id} does not exist')

    return TaskResponse(**task)


@task_api_router.get(path="/all",
                     responses={
                         status.HTTP_200_OK: {
                             'model': TaskListResponse
                         }
                     }
                     )
async def get_all_tasks(data: Annotated[TaskListRequest, Depends()]):
    tasks = fetch_tasks(page=data.page, size=data.size)
    if tasks is None:
        raise GenericNotFoundException(detail='This instance has no tasks')

    return TaskListResponse(tasks=tasks)

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
