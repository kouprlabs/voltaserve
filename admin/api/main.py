import logging
import time

import psycopg2
from fastapi import FastAPI, Request, Response, status

from routers import users_api_router
from dependencies import settings

logging.basicConfig(level=logging.DEBUG)

app = FastAPI(debug=True)

app.include_router(users_api_router)


@app.middleware("http")
async def add_process_time_header(request: Request, call_next):
    start_time = time.time()
    response = await call_next(request)
    response.headers["X-Process-Time-Ms"] = str(round((time.time() - start_time) * 1000, 4))
    return response


@app.get('/')
async def root():
    return {"message": "Hello, it is root of admin microservice!"}


@app.get('/liveness')
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
    except:
        return Response(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE
        )


@app.get('/readiness')
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
    except:
        return Response(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE
        )