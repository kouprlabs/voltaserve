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
  emailUpdateToken?: string
  emailUpdateValue?: string
  picture?: string
  createTime: string
  updateTime?: string
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
  picture?: string
  createTime?: string
  updateTime?: string
}

export type UpdateOptions = {
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
  emailUpdateToken?: string
  emailUpdateValue?: string
  picture?: string
  createTime?: string
  updateTime?: string
}

export interface UserRepo {
  findByID(id: string): Promise<User>
  findByUsername(username: string): Promise<User>
  findByEmail(email: string): Promise<User>
  findByRefreshTokenValue(refreshTokenValue: string): Promise<User>
  findByResetPasswordToken(resetPasswordToken: string): Promise<User>
  findByEmailConfirmationToken(emailConfirmationToken: string): Promise<User>
  findByEmailUpdateToken(emailUpdateToken: string): Promise<User>
  findByPicture(picture: string): Promise<User>
  isUsernameAvailable(username: string): Promise<boolean>
  insert(data: InsertOptions): Promise<User>
  update(data: UpdateOptions): Promise<User>
  delete(id: string): Promise<void>
}
