// Copyright 2024 Piotr Łoboda.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { Request } from 'express'

export interface UserIdPostRequest extends Request {
  id: string
}

export interface UserSuspendPostRequest extends UserIdPostRequest {
  suspend: boolean
}

export interface UserAdminPostRequest extends UserIdPostRequest {
  makeAdmin: boolean
}

export type Pagination = {
  page: string
  size: string
}

export type UserId = {
  id: string
}

export interface PaginatedRequest extends Request {
  query: Pagination
}

export interface UserIdRequest extends Request {
  query: UserId
}
