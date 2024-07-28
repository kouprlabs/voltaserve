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

import jwt
import psycopg2
from fastapi import FastAPI, Request, Response, status, Depends
from fastapi.middleware.cors import CORSMiddleware
from fastapi_utils.tasks import repeat_every

from .dependencies import settings, JWTBearer
from .routers import users_api_router, group_api_router, organization_api_router, workspace_api_router, \
    invitation_api_router, index_api_router

app = FastAPI(debug=True)

app.add_middleware(
    CORSMiddleware,
    allow_origins=settings.cors_origins.split(','),
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"]
)

logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s', level=logging.DEBUG)
logger = logging.getLogger('voltaserve-admin-api')

app.include_router(users_api_router)
app.include_router(group_api_router)
app.include_router(organization_api_router)
# app.include_router(task_api_router)
app.include_router(workspace_api_router)
app.include_router(invitation_api_router)
app.include_router(index_api_router)


@app.on_event("startup")
@repeat_every(seconds=settings.admin_token_expiration)
async def regenerate_admin_secret() -> None:
    jwt.api_jws.PyJWS.header_typ = False
    token = jwt.encode(payload={"sub": "SUPERUSER",
                                "iat": time.time(),
                                "iss": settings.url,
                                "aud": settings.url,
                                "exp": time.time() + settings.admin_token_expiration},
                       key=settings.jwt_secret,
                       algorithm=settings.jwt_algorithm
                       )

    logger.info(f"Admin token: \n{token}")
    settings.admin_token = token
    del token


@app.middleware("http")
async def add_process_time_header(request: Request, call_next):
    start_time = time.time()
    response = await call_next(request)
    response.headers["X-Process-Time-Ms"] = str(round((time.time() - start_time) * 1000, 4))
    return response


@app.get('/', tags=['main'], dependencies=[Depends(JWTBearer())])
async def root():
    return {"detail": "Hello, it is root of admin microservice!"}


@app.get('/liveness', tags=['liveness'])
async def liveness():
    try:
        psycopg2.connect(host=settings.db_host,
                         user=settings.db_user,
                         password=settings.db_password,
                         dbname=settings.db_name,
                         port=settings.db_port)
        return Response(
            status_code=status.HTTP_204_NO_CONTENT,
        )
    except Exception as e:
        return Response(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            content=e
        )


@app.get('/readiness', tags=['readiness'])
async def readiness():
    try:
        psycopg2.connect(host=settings.db_host,
                         user=settings.db_user,
                         password=settings.db_password,
                         dbname=settings.db_name,
                         port=settings.db_port)
        return Response(
            status_code=status.HTTP_204_NO_CONTENT,
        )
    except Exception as e:
        return Response(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE,
            content=e
        )
