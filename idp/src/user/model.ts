// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.

export enum Strategy {
  Local = 'local',
  Apple = 'apple',
}

export type User = {
  id: string
  fullName: string
  username: string
  email: string
  passwordHash: string
  refreshTokenValue?: string
  refreshTokenExpiry?: string
  resetPasswordToken?: string
  emailConfirmationToken?: string
  isEmailConfirmed: boolean
  isAdmin: boolean
  isActive: boolean
  emailUpdateToken?: string
  emailUpdateValue?: string
  picture?: string
  failedAttempts: number
  lockedUntil?: string
  strategy: Strategy
  createTime: string
  updateTime?: string
}

export function isLocalStrategy(user: User): boolean {
  return user.strategy === Strategy.Local
}

export function isAppleStrategy(user: User): boolean {
  return user.strategy === Strategy.Apple
}

export type InsertOptions = {
  id: string
  fullName?: string
  username?: string
  email?: string
  passwordHash?: string
  refreshTokenValue?: string
  refreshTokenExpiry?: string
  resetPasswordToken?: string
  emailConfirmationToken?: string
  isEmailConfirmed?: boolean
  isAdmin?: boolean
  isActive?: boolean
  picture?: string
  strategy: Strategy
  createTime?: string
  updateTime?: string
}

export type UpdateOptions = {
  id: string
  fullName?: string
  username?: string
  email?: string
  passwordHash?: string | null
  refreshTokenValue?: string | null
  refreshTokenExpiry?: string | null
  resetPasswordToken?: string | null
  emailConfirmationToken?: string | null
  isEmailConfirmed?: boolean
  emailUpdateToken?: string | null
  emailUpdateValue?: string | null
  picture?: string | null
  failedAttempts?: number
  lockedUntil?: string | null
  createTime?: string
  updateTime?: string | null
}

export interface UserRepo {
  findById(id: string): Promise<User>
  findByUsername(username: string): Promise<User>
  findByEmail(email: string): Promise<User>
  findByRefreshTokenValue(refreshTokenValue: string): Promise<User>
  findByResetPasswordToken(resetPasswordToken: string): Promise<User>
  findByEmailConfirmationToken(emailConfirmationToken: string): Promise<User>
  findByEmailUpdateToken(emailUpdateToken: string): Promise<User>
  findByPicture(picture: string): Promise<User>
  listAllPaginated(page: number, size: number): Promise<User>
  isUsernameAvailable(username: string): Promise<boolean>
  insert(data: InsertOptions): Promise<User>
  update(data: UpdateOptions): Promise<User>
  delete(id: string): Promise<void>
}
