# Copyright (c) 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

from typing import Annotated

from aiohttp import ClientSession
from fastapi import APIRouter, status, Depends

from ..dependencies import settings, JWTBearer
from ..log import base_logger
from ..errors import UnknownApiError, NotFoundError
from ..models import VersionRequest, VersionResponse

overview_api_router = APIRouter(
    prefix="/overview", tags=["overview"], dependencies=[Depends(JWTBearer())]
)

logger = base_logger.getChild("overview")


async def get_dockerhub_version(
    sess: ClientSession, id: str, response: dict, params: dict
) -> dict:
    try:
        async with sess.get(
            f"https://hub.docker.com/v2/repositories/voltaserve/{id}/tags",
            params=params,
        ) as resp:
            if resp.status == 200:
                tags = await resp.json()
                latest_digest = next(
                    tag["digest"] for tag in tags["results"] if tag["name"] == "latest"
                )
                response["latestVersion"] = max(
                    tag["name"]
                    for tag in tags["results"]
                    if tag["digest"] == latest_digest and tag["name"] != "latest"
                )
                response["location"] = (
                    f"https://hub.docker.com/layers/voltaserve/{id}"
                    f"/{response['latestVersion']}/images/{latest_digest}"
                )
            else:
                response["latestVersion"] = "UNKNOWN"
                response["location"] = f"https://hub.docker.com/layers/voltaserve/{id}"
    except Exception as e:
        logger.exception(e)
        response["latestVersion"] = "UNKNOWN"
        response["location"] = f"https://hub.docker.com/layers/voltaserve/{id}"

    return response


async def get_local_version(sess: ClientSession, url: str, response: dict) -> dict:
    try:
        async with sess.get(f"{url}/version") as local_resp:
            if local_resp.status == 200:
                response["currentVersion"] = (await local_resp.json())["version"]
            else:
                response["currentVersion"] = "UNKNOWN"
    except Exception as e:
        logger.exception(e)
        response["currentVersion"] = "UNKNOWN"

    return response


@overview_api_router.get(
    path="/version/internal", responses={status.HTTP_200_OK: {"model": VersionResponse}}
)
async def get_internal_version(data: Annotated[VersionRequest, Depends()]):
    if data.id not in (
        "api",
        "conversion",
        "idp",
        "language",
        "mosaic",
        "ui",
        "webdav",
        "console",
    ):
        return NotFoundError(message=f"Microservice {data.id} not found")

    try:
        urls = settings.model_dump()
        response = {"name": data.id}
        params = {"page_size": 50, "page": 1, "ordering": "last_updated", "name": ""}
        async with ClientSession() as sess:
            response = await get_dockerhub_version(sess, data.id, response, params)
            if data.id == "console":
                response["currentVersion"] = "3.0.0"
            elif data.id == "ui":
                response["currentVersion"] = ""
            else:
                response = await get_local_version(
                    sess, urls[f"{data.id.upper()}_URL"], response
                )

            response["updateAvailable"] = (
                response["latestVersion"] > response["currentVersion"]
                if response["currentVersion"] != ""
                and response["latestVersion"] != "UNKNOWN"
                else None
            )

        return VersionResponse(**response)
    except Exception as e:
        logger.exception(e)
        return UnknownApiError()
