import uvicorn

from .dependencies import settings

if __name__ == "__main__":
    log_config = uvicorn.config.LOGGING_CONFIG
    log_config["formatters"]["access"]["fmt"] = ('%(asctime)s - %(name)s - %(levelname)s - %(client_addr)s - '
                                                 '"%(request_line)s" %(status_code)s')
    log_config["formatters"]["default"]["fmt"] = '%(asctime)s - %(name)s - %(levelname)s - %(message)s'

    uvicorn.run(app="api.main:app",
                host=settings.host,
                port=settings.port,
                reload=False,
                workers=settings.workers,
                log_config=log_config)
