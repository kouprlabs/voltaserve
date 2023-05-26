import fs from 'fs/promises'
import { ErrorCode, newError } from '@/infra/error'
import { hashPassword, verifyPassword } from '@/infra/password'
import search, { USER_SEARCH_INDEX } from '@/infra/search'
import userRepo, { User } from '@/user/repo'

export type UserDTO = {
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

export async function getUser(id: string): Promise<UserDTO> {
  return mapEntity(await userRepo.findByID(id))
}

export async function getByPicture(picture: string): Promise<UserDTO> {
  return mapEntity(await userRepo.findByPicture(picture))
}

export async function updateFullName(
  id: string,
  options: UserUpdateFullNameOptions
): Promise<UserDTO> {
  let user = await userRepo.findByID(id)
  user = await userRepo.update({ id: user.id, fullName: options.fullName })
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
): Promise<UserDTO> {
  let user = await userRepo.findByID(id)
  user = await userRepo.update({
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
): Promise<UserDTO> {
  let user = await userRepo.findByID(id)
  if (verifyPassword(options.currentPassword, user.passwordHash)) {
    user = await userRepo.update({
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
): Promise<UserDTO> {
  const picture = await fs.readFile(path, { encoding: 'base64' })
  const { id: userId } = await userRepo.findByID(id)
  const user = await userRepo.update({
    id: userId,
    picture: `data:${contentType};base64,${picture}`,
  })
  return mapEntity(user)
}

export async function deletePicture(id: string): Promise<UserDTO> {
  let user = await userRepo.findByID(id)
  user = await userRepo.update({ id: user.id, picture: null })
  return mapEntity(user)
}

export async function deleteUser(id: string, options: UserDeleteOptions) {
  const user = await userRepo.findByID(id)
  if (verifyPassword(options.password, user.passwordHash)) {
    await userRepo.delete(user.id)
    await search.index(USER_SEARCH_INDEX).deleteDocuments([user.id])
  } else {
    throw newError({ code: ErrorCode.InvalidPassword })
  }
}

export function mapEntity(entity: User): UserDTO {
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
