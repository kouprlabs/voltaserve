import time
import logging

import psycopg2
from fastapi import FastAPI, Request, Response, status
from routers import users_api_router

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
        psycopg2.connect(host='192.168.1.254', user='root', dbname='voltaserve', port=26257)
        return Response(
            status_code=status.HTTP_200_OK,
        )
    except:
        return Response(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE
        )


@app.get('/readiness')
async def readiness():
    try:
        psycopg2.connect(host='192.168.1.254', user='root', dbname='voltaserve', port=26257)
        return Response(
            status_code=status.HTTP_200_OK,
        )
    except:
        return Response(
            status_code=status.HTTP_503_SERVICE_UNAVAILABLE
        )
