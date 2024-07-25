import logging
import time

import psycopg2
from fastapi import FastAPI, Request, Response, status

from .dependencies import settings
from .routers import users_api_router, group_api_router, organization_api_router, task_api_router, \
    workspace_api_router, invitation_api_router, token_api_router

app = FastAPI(debug=True)

logging.basicConfig(format='%(asctime)s - %(name)s - %(levelname)s - %(message)s', level=logging.DEBUG)
logger = logging.getLogger(__name__)

app.include_router(users_api_router)
app.include_router(group_api_router)
app.include_router(organization_api_router)
app.include_router(task_api_router)
app.include_router(workspace_api_router)
app.include_router(invitation_api_router)
app.include_router(token_api_router)


@app.middleware("http")
async def add_process_time_header(request: Request, call_next):
    start_time = time.time()
    response = await call_next(request)
    response.headers["X-Process-Time-Ms"] = str(round((time.time() - start_time) * 1000, 4))
    return response


@app.get('/', tags=['main'])
async def root():
    return {"message": "Hello, it is root of admin microservice!"}


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
