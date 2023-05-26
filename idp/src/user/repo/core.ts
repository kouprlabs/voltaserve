export type User = {
  id: string
  fullName: string
  username: string
  email: string
  passwordHash: string
  refreshTokenValue?: string
  refreshTokenValidTo?: number
  resetPasswordToken?: string
  emailConfirmationToken?: string
  isEmailConfirmed: boolean
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
  refreshTokenValidTo?: number
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
  refreshTokenValidTo?: number
  resetPasswordToken?: string
  emailConfirmationToken?: string
  isEmailConfirmed?: boolean
  picture?: string
  createTime?: string
  updateTime?: string
}

export type Field =
  | 'id'
  | 'full_name'
  | 'username'
  | 'email'
  | 'password_hash'
  | 'refresh_token_value'
  | 'refresh_token_valid_to'
  | 'reset_password_token'
  | 'email_confirmation_token'
  | 'is_email_confirmed'
  | 'picture'

export interface UserRepo {
  find(field: Field, value: any, canThrow?: boolean): Promise<User>
  findByPicture(picture: string): Promise<User>
  insert(data: InsertOptions): Promise<User>
  update(data: UpdateOptions): Promise<User>
  delete(id: string): Promise<void>
}
