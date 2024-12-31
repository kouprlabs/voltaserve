# Copyright (c) 2023 Anass Bouassaba.
#
# Use of this software is governed by the Business Source License
# included in the file LICENSE in the root of this repository.
#
# As of the Change Date specified in that file, in accordance with
# the Business Source License, use of this software will be governed
# by the GNU Affero General Public License v3.0 only, included in the file
# AGPL-3.0-only in the root of this repository.

from fastapi import APIRouter, Depends, Response, status

from api.database import grant_user_permission, revoke_user_permission
from api.dependencies import JWTBearer, redis_conn
from api.models.userpermission import (
    UserPermissionGrantRequest,
    UserPermissionRevokeRequest,
)

user_permission_api_router = APIRouter(
    prefix="/user_permission", tags=["permission"], dependencies=[Depends(JWTBearer())]
)


@user_permission_api_router.post(path="/grant", status_code=status.HTTP_200_OK)
async def post_grant_user_permission(data: UserPermissionGrantRequest):
    await redis_conn.delete(f"{data.resourceType}:{data.resourceId}")
    grant_user_permission(
        user_id=data.userId,
        resource_id=data.resourceId,
        permission=data.permission,
    )
    return Response(status_code=status.HTTP_200_OK)


@user_permission_api_router.post(path="/revoke", status_code=status.HTTP_200_OK)
async def post_revoke_user_permission(data: UserPermissionRevokeRequest):
    await redis_conn.delete(f"{data.resourceType}:{data.resourceId}")
    revoke_user_permission(
        user_id=data.userId,
        resource_id=data.resourceId,
    )
    return Response(status_code=status.HTTP_200_OK)
