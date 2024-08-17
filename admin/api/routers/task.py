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

from ..database import fetch_task, fetch_tasks
from ..errors import NotFoundError, EmptyDataException, NoContentError, NotFoundException, \
    UnknownApiError
from ..models import TaskResponse, TaskRequest, TaskListResponse, TaskListRequest
from ..tasks.celery_app import celery_app

task_api_router = APIRouter(
    prefix='/task',
    tags=['task'],
)

user_task_api_router = APIRouter(
    prefix='/user',
    tags=['task'],
)

admin_task_api_router = APIRouter(
    prefix='/admin',
    tags=['task'],
)


# --- GET --- #
@user_task_api_router.get(path="",
                          responses={
                              status.HTTP_200_OK: {
                                  'model': TaskResponse
                              }}
                          )
async def get_user_task(data: Annotated[TaskRequest, Depends()]):
    try:
        task = fetch_task(_id=data.id)
        if task is None:
            return NotFoundError(message=f'Task with id={data.id} does not exist')

        return TaskResponse(**task)
    except Exception as e:
        return UnknownApiError(message=str(e))


@user_task_api_router.get(path="/all",
                          responses={
                              status.HTTP_200_OK: {
                                  'model': TaskListResponse
                              }
                          }
                          )
async def get_all_user_tasks(data: Annotated[TaskListRequest, Depends()]):
    try:
        tasks, count = fetch_tasks(page=data.page, size=data.size)

        return TaskListResponse(data=tasks, totalElements=count, page=data.page, size=data.size)
    except EmptyDataException:
        return NoContentError()
    except NotFoundException as e:
        return NotFoundError(message=str(e))
    except Exception as e:
        return UnknownApiError(message=str(e))


@admin_task_api_router.get(path="/all",
                           # responses={
                           #     status.HTTP_200_OK: {
                           #         'model': TaskListResponse
                           #     }
                           # }
                           )
async def get_all_admin_tasks():
    try:
        return celery_app.control.inspect().reserved()
    except Exception as e:
        return UnknownApiError(message=str(e))


@admin_task_api_router.get(path="/active",
                           # responses={
                           #     status.HTTP_200_OK: {
                           #         'model': TaskListResponse
                           #     }
                           # }
                           )
async def get_active_admin_tasks():
    try:
        return celery_app.control.inspect().active()
    except Exception as e:
        return UnknownApiError(message=str(e))


@admin_task_api_router.get(path="/scheduled",
                           # responses={
                           #     status.HTTP_200_OK: {
                           #         'model': TaskListResponse
                           #     }
                           # }
                           )
async def get_scheduled_admin_tasks():
    try:
        return celery_app.control.inspect().scheduled()
    except Exception as e:
        return UnknownApiError(message=str(e))

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
