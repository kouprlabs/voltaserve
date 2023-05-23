import fs from 'fs/promises'
import { UserEntity, UserRepo } from '@/infra/db'
import { ErrorCode, newError } from '@/infra/error'
import { hashPassword, verifyPassword } from '@/infra/password'
import search, { USER_SEARCH_INDEX } from '@/infra/search'

export type User = {
  id: string
  fullName: string
  picture: string
  email: string
  username: string
}

export type UserUpdateFullNameOptions = {
  fullName: string
}

export type UserUpdateEmailOptions = {
  email: string
}

export type UserUpdatePasswordOptions = {
  currentPassword: string
  newPassword: string
}

export type UserDeleteOptions = {
  password: string
}

export async function getUser(id: string): Promise<User> {
  return mapEntity(await UserRepo.find('id', id, true))
}

export async function getByPicture(picture: string): Promise<User> {
  return mapEntity(await UserRepo.findByPicture(picture))
}

export async function updateFullName(
  id: string,
  options: UserUpdateFullNameOptions
): Promise<User> {
  let user = await UserRepo.find('id', id, true)
  user = await UserRepo.update({ id: user.id, fullName: options.fullName })
  await search.index(USER_SEARCH_INDEX).updateDocuments([
    {
      ...user,
      fullName: user.fullName,
    },
  ])
  return mapEntity(user)
}

export async function updateEmail(
  id: string,
  options: UserUpdateEmailOptions
): Promise<User> {
  let user = await UserRepo.find('id', id, true)
  user = await UserRepo.update({
    id: user.id,
    email: options.email,
    username: options.email,
  })
  await search.index(USER_SEARCH_INDEX).updateDocuments([
    {
      ...user,
      email: user.email,
      username: user.email,
    },
  ])
  return mapEntity(user)
}

export async function updatePassword(
  id: string,
  options: UserUpdatePasswordOptions
): Promise<User> {
  let user = await UserRepo.find('id', id, true)
  if (verifyPassword(options.currentPassword, user.passwordHash)) {
    user = await UserRepo.update({
      id: user.id,
      passwordHash: hashPassword(options.newPassword),
    })
    return mapEntity(user)
  } else {
    throw newError({ code: ErrorCode.PasswordValidationFailed })
  }
}

export async function updatePicture(
  id: string,
  path: string,
  contentType: string
): Promise<User> {
  const picture = await fs.readFile(path, { encoding: 'base64' })
  const { id: userId } = await UserRepo.find('id', id, true)
  const user = await UserRepo.update({
    id: userId,
    picture: `data:${contentType};base64,${picture}`,
  })
  return mapEntity(user)
}

export async function deletePicture(id: string): Promise<User> {
  let user = await UserRepo.find('id', id, true)
  user = await UserRepo.update({ id: user.id, picture: null })
  return mapEntity(user)
}

export async function deleteUser(id: string, options: UserDeleteOptions) {
  const user = await UserRepo.find('id', id, true)
  if (verifyPassword(options.password, user.passwordHash)) {
    await UserRepo.delete(user.id)
    await search.index(USER_SEARCH_INDEX).deleteDocuments([user.id])
  } else {
    throw newError({ code: ErrorCode.InvalidPassword })
  }
}

export function mapEntity(entity: UserEntity): User {
  const user = {
    id: entity.id,
    email: entity.email,
    username: entity.username,
    fullName: entity.fullName,
    picture: entity.picture,
  }
  Object.keys(user).forEach(
    (index) => !user[index] && user[index] !== undefined && delete user[index]
  )
  return user
}
