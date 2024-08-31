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
import uuid

from fastapi import FastAPI, Request, Response, status, HTTPException
from fastapi.middleware.cors import CORSMiddleware

from .dependencies import settings, conn
from .errors import ServiceUnavailableError, ForbiddenError
from .models import GenericServiceUnavailableResponse, GenericErrorResponse
from .routers import group_api_router, organization_api_router, workspace_api_router, \
    invitation_api_router, index_api_router, users_api_router, overview_api_router

if settings.LOG_FORMAT == 'JSON':
    req_fmt = ('{"timestamp":"%(asctime)s",'
               '"logger_name":"%(name)s",'
               '"log_level":"%(levelname)s",'
               '"message":"%(message)s",'
               '"type":"%(type)s",'
               '"identifier":"%(identifier)s",'
               '"path":"%(path)s",'
               '"method":"%(method)s",'
               '"headers":"%(headers)s",'
               '"query_params":"%(query_params)s",'
               '"path_params":"%(path_params)s"}')
    resp_fmt = ('{"timestamp":"%(asctime)s",'
                '"logger_name":"%(name)s",'
                '"log_level":"%(levelname)s",'
                '"message":"%(message)s",'
                '"type":"%(type)s",'
                '"identifier":"%(identifier)s",'
                '"code":"%(code)s",'
                '"headers":"%(headers)s"'
                )
    base_fmt = ('{"timestamp":"%(asctime)s",'
                '"logger_name":"%(name)s",'
                '"log_level":"%(levelname)s",'
                '"message":"%(message)s"}'
                )
elif settings.LOG_FORMAT == "PLAIN":
    req_fmt = ('%(asctime)s|%(name)s|%(levelname)s|%(type)s|%(identifier)s|%(path)s|'
               '%(method)s|%(headers)s|%(query_params)s|%(path_params)s|%(message)s')
    resp_fmt = '%(asctime)s|%(name)s|%(levelname)s|%(type)s|%(identifier)s|%(code)s|%(headers)s|%(message)s'
    base_fmt = '%(asctime)s|%(name)s|%(levelname)s|%(message)s'
else:
    raise ValueError('Wrong logging format, available JSON and PLAIN')

logging.basicConfig(format=base_fmt, level=settings.LOG_LEVEL)
logger = logging.getLogger('voltaserve.console.api')


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

if settings.LOG_LEVEL == 'DEBUG':
    req_handler = logging.StreamHandler()
    req_handler.setFormatter(logging.Formatter(req_fmt))

    resp_handler = logging.StreamHandler()
    resp_handler.setFormatter(logging.Formatter(resp_fmt))

    req_logger = logger.getChild('requests')
    req_logger.addHandler(req_handler)

    resp_logger = logger.getChild('responses')
    resp_logger.addHandler(resp_handler)


    @app.middleware("http")
    async def log_requests(request: Request, call_next):
        identifier = uuid.uuid4()
        x = {'type': 'request',
             'identifier': identifier,
             'path': request.url.path,
             'method': request.method,
             'headers': dict(request.headers),
             'query_params': dict(request.query_params),
             'path_params': dict(request.path_params)
             }
        req_logger.debug(msg='', extra=x)

        response: Response = await call_next(request)
        resp_logger.debug(msg='', extra={'type': 'response',
                                         'identifier': identifier,
                                         'code': response.status_code,
                                         'headers': dict(response.headers)
                                         })

        return response

app.include_router(users_api_router)
app.include_router(group_api_router)
app.include_router(organization_api_router)
app.include_router(workspace_api_router)
app.include_router(invitation_api_router)
app.include_router(index_api_router)
app.include_router(overview_api_router)


@app.middleware("http")
async def add_process_time_header(request: Request, call_next):
    start_time = time.time()
    response = await call_next(request)
    response.headers["X-Process-Time-Ms"] = str(round((time.time() - start_time) * 1000, 4))
    return response


@app.get('/', tags=['main'])
async def root():
    return {"detail": "Hello, it is root of console microservice!"}


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
