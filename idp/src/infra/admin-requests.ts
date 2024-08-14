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

export type Pagination = {
  page: string
  size: string
}

export interface UserIdRequest extends Request {
  id: string
}

export interface UserSuspendRequest extends UserIdRequest {
  suspend: boolean
}

export type UserCreationDate = {
  userCreationDate: string
}

export type UserUpdateDate = {
  id: string
}

export interface PaginatedRequest extends Request {
  query: Pagination
}
