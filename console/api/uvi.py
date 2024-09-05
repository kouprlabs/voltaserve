# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

import uvicorn

from .dependencies import settings

if __name__ == "__main__":
    log_config = uvicorn.config.LOGGING_CONFIG
    if settings.LOG_FORMAT == "JSON":
        log_config["formatters"]["access"]["fmt"] = (
            '{"timestamp":"%(asctime)s",'
            '"logger_name":"%(name)s",'
            '"log_level":"%(levelname)s",'
            '"source_address":"%(client_addr)s",'
            '"path":"%(request_line)s",'
            '"status_code":"%(status_code)s"}'
        )
        log_config["formatters"]["default"]["fmt"] = (
            '{"timestamp":"%(asctime)s",'
            '"logger_name":"%(name)s",'
            '"log_level":"%(levelname)s",'
            '"message":"%(message)s"}'
        )
    elif settings.LOG_FORMAT == "PLAIN":
        log_config["formatters"]["access"]["fmt"] = (
            "%(asctime)s|%(name)s|%(levelname)s|%(client_addr)s|"
            "%(request_line)s|%(status_code)s"
        )
        log_config["formatters"]["default"][
            "fmt"
        ] = "%(asctime)s|%(name)s|%(levelname)s|%(message)s"
    else:
        raise ValueError("Wrong logging format, available JSON and PLAIN")
    if settings.LOG_LEVEL == "DEBUG":
        log_config["loggers"]["uvicorn.access"]["level"] = "CRITICAL"

    uvicorn.run(
        app="api.main:app",
        host=settings.HOST,
        port=settings.PORT,
        reload=False,
        workers=settings.WORKERS,
        log_config=log_config,
    )
