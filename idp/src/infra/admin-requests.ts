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
import { User } from '@/user/model'

export interface UserIdPostRequest extends Request {
  id: string
}

export type SearchRequest = {
  page: string
  size: string
  query: string
}

export interface UserSuspendPostRequest extends UserIdPostRequest {
  suspend: boolean
}

export interface UserAdminPostRequest extends UserIdPostRequest {
  makeAdmin: boolean
}

export type UserId = {
  id: string
}

export interface SearchPaginatedRequest extends Request {
  query: SearchRequest
}

export interface UserIdRequest extends Request {
  query: UserId
}

export interface UserSearchResponse {
  data: User[]
  page: number
  size: number
  totalElements: number
}

export interface UserUpdateAdminRequest {
  fullName: string
  email: string
}
