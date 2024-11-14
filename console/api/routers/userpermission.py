# Copyright 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file licenses/BSL.txt.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# licenses/AGPL.txt.

from typing import Annotated

from fastapi import APIRouter, Depends, status

from api.database import grant_user_permission, revoke_user_permission
from api.dependencies import JWTBearer
from api.models.userpermission import (
    UserPermissionGrantRequest,
    UserPermissionRevokeRequest,
)

userpermission_api_router = APIRouter(
    prefix="/userpermission", tags=["permission"], dependencies=[Depends(JWTBearer())]
)


@userpermission_api_router.post(path="/grant", responses={status.HTTP_200_OK: {}})
async def post_grant_user_permission(
    data: Annotated[UserPermissionGrantRequest, Depends()]
):
    grant_user_permission(
        user_id=data.user_id,
        resource_id=data.resource_id,
        permission=data.permission,
    )


@userpermission_api_router.post(path="/revoke", responses={status.HTTP_200_OK: {}})
async def post_revoke_user_permission(
    data: Annotated[UserPermissionRevokeRequest, Depends()]
):
    revoke_user_permission(
        user_id=data.user_id,
        resource_id=data.resource_id,
    )
