// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { z, ZodError } from 'zod'
import { ErrorCode, ErrorData, newError } from '@/infra/error/core.ts'
import { getConfig } from '@/config/config.ts'
import { Buffer } from 'node:buffer'

export function parseValidationError(result: ZodError): ErrorData {
  let message: string | undefined
  let userMessage: string | undefined
  if (result.errors) {
    message = result.errors
      .map((e) => e.message ? `${e.message} (${e.path.join('.')}).` : undefined)
      .join(' ')
    userMessage = result.errors
      .map((e) => e.message ? `${e.message} (${e.path.join('.')}).` : undefined)
      .join(' ')
  }
  return newError({
    code: ErrorCode.RequestValidationError,
    message,
    userMessage,
  })
}

type ZodValidationResult<T> = {
  success: true
  data: T
} | {
  success: false
  error: ZodError
  data: T
}

export function handleValidationError<T>(result: ZodValidationResult<T>) {
  if (!result.success) {
    throw parseValidationError(result.error)
  }
}

export const password = z
  .string()
  .min(
    getConfig().password.minLength,
    `Password must be at least ${getConfig().password.minLength} characters long`,
  )
  .regex(
    new RegExp(`(?=(.*[a-z]){${getConfig().password.minLowercase},})`),
    `Password must contain at least ${getConfig().password.minLowercase} lowercase character(s)`,
  )
  .regex(
    new RegExp(`(?=(.*[A-Z]){${getConfig().password.minUppercase},})`),
    `Password must contain at least ${getConfig().password.minUppercase} uppercase character(s)`,
  )
  .regex(
    new RegExp(`(?=(.*[0-9]){${getConfig().password.minNumbers},})`),
    `Password must contain at least ${getConfig().password.minNumbers} number(s)`,
  )
  .regex(
    new RegExp(
      `(?=(.*[!@#$%^&*()_+\\-=\\[\\]{};':"\\|,.<>\\/?]){${getConfig().password.minSymbols},})`,
    ),
    `Password must contain at least ${getConfig().password.minSymbols} symbol(s)`,
  )

export const picture = z
  .string()
  .optional()
  .refine((value) => {
    if (!value) {
      return true
    }
    try {
      return Buffer.from(value, 'base64').length <= 3000000
    } catch {
      return false
    }
  }, { message: 'Picture must be a valid Base64 string and <= 3MB' })

export const email = z.string().email().trim().max(255)

export const fullName = z.string().nonempty().trim().max(255)

export const token = z.string().nonempty()
