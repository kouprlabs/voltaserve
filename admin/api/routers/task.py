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
from ..exceptions import GenericNotFoundException
from ..models import GenericNotFoundResponse, TaskResponse, TaskRequest, TaskListResponse, TaskListRequest
from ..tasks.celery_app import celery_app

task_api_router = APIRouter(
    prefix='/task',
    tags=['task'],
    responses={
        status.HTTP_404_NOT_FOUND: {
            'model': GenericNotFoundResponse
        }
    }
)

user_task_api_router = APIRouter(
    prefix='/user',
    tags=['task'],
    responses={
        status.HTTP_404_NOT_FOUND: {
            'model': GenericNotFoundResponse
        }
    }
)

admin_task_api_router = APIRouter(
    prefix='/admin',
    tags=['task'],
    responses={
        status.HTTP_404_NOT_FOUND: {
            'model': GenericNotFoundResponse
        }
    }
)


# --- GET --- #
@user_task_api_router.get(path="",
                          responses={
                              status.HTTP_200_OK: {
                                  'model': TaskResponse
                              }}
                          )
async def get_user_task(data: Annotated[TaskRequest, Depends()]):
    task = fetch_task(_id=data.id)
    if task is None:
        raise GenericNotFoundException(detail=f'Task with id={data.id} does not exist')

    return TaskResponse(**task)


@user_task_api_router.get(path="/all",
                          responses={
                              status.HTTP_200_OK: {
                                  'model': TaskListResponse
                              }
                          }
                          )
async def get_all_user_tasks(data: Annotated[TaskListRequest, Depends()]):
    tasks, count = fetch_tasks(page=data.page, size=data.size)
    if tasks is None:
        raise GenericNotFoundException(detail='This instance has no tasks')

    return TaskListResponse(data=tasks, totalElements=count['count'], page=data.page, size=data.size)


@admin_task_api_router.get(path="/all",
                           # responses={
                           #     status.HTTP_200_OK: {
                           #         'model': TaskListResponse
                           #     }
                           # }
                           )
async def get_all_admin_tasks():
    return celery_app.control.inspect().reserved()


@admin_task_api_router.get(path="/active",
                           # responses={
                           #     status.HTTP_200_OK: {
                           #         'model': TaskListResponse
                           #     }
                           # }
                           )
async def get_active_admin_tasks():
    return celery_app.control.inspect().active()


@admin_task_api_router.get(path="/scheduled",
                           # responses={
                           #     status.HTTP_200_OK: {
                           #         'model': TaskListResponse
                           #     }
                           # }
                           )
async def get_scheduled_admin_tasks():
    return celery_app.control.inspect().scheduled()

# --- PATCH --- #

# --- POST --- #

# --- PUT --- #

# --- DELETE --- #
