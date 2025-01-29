// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { getConfig } from '@/config/config'
import store from '@/store/configure-store'
import { errorOccurred } from '@/store/ui/error'
import { errorToString } from './error'
import { getAccessToken, getAccessTokenOrRedirect } from './token'

export type FetcherOptions = {
  url: string
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH' | 'OPTIONS' | 'HEAD'
  body?: BodyInit | null
  contentType?: string
  redirect?: boolean
  authenticate?: boolean
  showError?: boolean
}

export function apiFetcher<T>(options: FetcherOptions) {
  return fetcher<T>({ ...options, url: `${getConfig().apiURL}${options.url}` })
}

export function idpFetcher<T>(options: FetcherOptions) {
  return fetcher<T>({
    ...options,
    url: `${getConfig().idpURL}${options.url}`,
  })
}

export function consoleFetcher<T>(options: FetcherOptions) {
  return fetcher<T>({
    ...options,
    url: `${getConfig().consoleURL}${options.url}`,
  })
}

export async function fetcher<T>({
  url,
  method,
  body,
  contentType,
  redirect,
  authenticate = true,
  showError = true,
}: FetcherOptions): Promise<T | undefined> {
  const headers: HeadersInit = {}
  if (!contentType) {
    headers['Content-Type'] = 'application/json'
  }
  if (authenticate) {
    headers['Authorization'] = `Bearer ${
      redirect ? getAccessTokenOrRedirect() : getAccessToken()
    }`
  }
  const response = await baseFetcher(
    url,
    {
      method,
      body,
      headers,
      credentials: authenticate ? 'include' : undefined,
    },
    redirect,
    showError,
  )
  try {
    if (response) {
      return (await response.json()) as T
    }
  } catch {
    // Ignored
  }
}

export async function baseFetcher(
  url: string,
  init: RequestInit,
  redirect = true,
  showError = true,
) {
  try {
    const response = await fetch(url, init)
    return handleResponse(response, redirect, showError)
  } catch {
    if (showError) {
      const message = 'Unexpected error occurred.'
      store.dispatch(errorOccurred(message))
      throw new Error(message)
    }
  }
}

async function handleResponse(
  response: Response,
  redirect = true,
  showError = true,
) {
  if (response.status <= 299) {
    return response
  } else {
    if (response.status === 401 && redirect) {
      window.location.href = '/sign-in'
    }
    let error
    try {
      error = await response.json()
    } catch {
      error = 'Oops! something went wrong.'
    }
    if (showError) {
      store.dispatch(errorOccurred(errorToString(error)))
    }
    throw error
  }
}
