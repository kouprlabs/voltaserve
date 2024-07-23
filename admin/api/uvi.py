import uvicorn

from dependencies import settings

if __name__ == "__main__":
    uvicorn.run("main:app",
                host=settings.host,
                port=settings.port,
                reload=True,
                workers=settings.workers)
