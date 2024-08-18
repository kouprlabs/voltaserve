# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

import logging
import time

from fastapi import FastAPI, Request, Response, status, HTTPException
from fastapi.middleware.cors import CORSMiddleware

from .dependencies import settings, conn
from .errors import ServiceUnavailableError, ForbiddenError
from .models import GenericServiceUnavailableResponse, GenericErrorResponse
from .routers import group_api_router, organization_api_router, workspace_api_router, \
    invitation_api_router, index_api_router, users_api_router, overview_api_router


async def custom_http_exception_handler(request: Request, exc: HTTPException):
    return ForbiddenError(
        status_code=exc.status_code,
        message=exc.detail
    )


app = FastAPI(root_path='/v1',
              debug=True,
              responses={
                  status.HTTP_204_NO_CONTENT: {
                      'model': None
                  },
                  status.HTTP_403_FORBIDDEN: {
                      'model': GenericErrorResponse
                  },
                  status.HTTP_404_NOT_FOUND: {
                      'model': GenericErrorResponse
                  },
                  status.HTTP_500_INTERNAL_SERVER_ERROR: {
                      'model': GenericErrorResponse
                  },
              },
              exception_handlers={
                  status.HTTP_403_FORBIDDEN: custom_http_exception_handler
              })

app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.SECURITY_CORS_ORIGINS.split(','),
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"]
)

logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s', level=logging.DEBUG)
logger = logging.getLogger('voltaserve-admin-api')

app.include_router(users_api_router)
app.include_router(group_api_router)
app.include_router(organization_api_router)
app.include_router(workspace_api_router)
app.include_router(invitation_api_router)
# app.include_router(index_api_router)
app.include_router(overview_api_router)


@app.middleware("http")
async def add_process_time_header(request: Request, call_next):
    start_time = time.time()
    response = await call_next(request)
    response.headers["X-Process-Time-Ms"] = str(round((time.time() - start_time) * 1000, 4))
    return response


@app.get('/', tags=['main'])
async def root():
    return {"detail": "Hello, it is root of admin microservice!"}


@app.get(path='/liveness',
         tags=['liveness'],
         status_code=status.HTTP_204_NO_CONTENT,
         responses={
             status.HTTP_204_NO_CONTENT: {},
             status.HTTP_503_SERVICE_UNAVAILABLE: {
                 'model': GenericServiceUnavailableResponse
             }
         }
         )
async def liveness(response: Response):
    try:
        with conn.cursor() as curs:
            curs.execute('SELECT 1;')

        response.status_code = status.HTTP_204_NO_CONTENT
        return None
    except Exception:
        return ServiceUnavailableError()


@app.get(path='/readiness',
         tags=['readiness'],
         status_code=status.HTTP_204_NO_CONTENT,
         responses={
             status.HTTP_204_NO_CONTENT: {},
             status.HTTP_503_SERVICE_UNAVAILABLE: {
                 'model': GenericServiceUnavailableResponse
             }
         }
         )
async def readiness(response: Response):
    try:
        with conn.cursor() as curs:
            curs.execute('SELECT 1;')

        response.status_code = status.HTTP_204_NO_CONTENT
        return None
    except Exception:
        return ServiceUnavailableError()
