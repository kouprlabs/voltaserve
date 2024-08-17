# Copyright 2024 Piotr Łoboda.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

from fastapi import APIRouter, status, Depends
from aiohttp import ClientSession

from ..database import fetch_version
from ..dependencies import JWTBearer
from ..errors import UnknownApiError
from ..models import VersionsResponse

overview_api_router = APIRouter(
    prefix='/overview',
    tags=['overview'],
    dependencies=[Depends(JWTBearer())]
)


@overview_api_router.get(path="/versions",
                         responses={
                             status.HTTP_200_OK: {
                                 'model': VersionsResponse
                             }
                         }
                         )
async def get_versions():
    try:
        loc = ['admin', 'database']
        ext = ['api', 'conversion', 'idp', 'language', 'mosaic', 'ui', 'webdav']
        data = []
        # series of requests to obtain versions
        async with ClientSession() as sess:
            for svc in ext:
                url = f"https://hub.docker.com/v2/repositories/voltaserve/{svc}/tags"
                params = {"page_size": 50, "page": 1, "ordering": "last_updated", "name": ""}

                async with sess.get(url, params=params) as resp:
                    tags = await resp.json()

                    latest_digest = next(l['digest'] for l in tags['results'] if l['name'] == 'latest')
                    latest_version = max(
                        l['name'] for l in tags['results'] if l['digest'] == latest_digest and l['name'] != 'latest')

                    data.append({
                        'name': svc,
                        'currentVersion': '2.0.1',
                        'latestVersion': latest_version,
                        'updateAvailable': latest_version > '2.0.1',
                        'location': f"https://hub.docker.com/layers/voltaserve/{svc}/"
                                    f"{latest_version}/images/{latest_digest}"
                    })

        data.extend([{
            'name': 'database',
            'currentVersion': fetch_version(),
            'latestVersion': '16.1',
            'updateAvailable': '16.1' > fetch_version(),
            'location': 'https://hub.docker.com/_/postgres'
        },
            {
                'name': 'admin',
                'currentVersion': '2.0.1',
                'latestVersion': '2.1.0',
                'updateAvailable': '2.1.0' > '2.0.1',
                'location': f"https://hub.docker.com/layers/voltaserve/admin/"
                            f"{latest_version}/images/{latest_digest}"
        }
        ])

        return VersionsResponse(data=data, totalElements=len(data), page=1, size=len(data))
    except Exception as e:
        return UnknownApiError(message=str(e))
