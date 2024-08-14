// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { getConfig } from '@/config/config'
import { getAccessToken, getAccessTokenOrRedirect } from '@/infra/token'
import store from '@/store/configure-store'
import { errorOccurred } from '@/store/ui/error'
import { errorToString } from './error'

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

export function adminFetcher<T>(options: FetcherOptions) {
  return fetcher<T>({
    ...options,
    url: `${getConfig().adminURL}${options.url}`,
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
  headers['Access-Control-Allow-Origin'] = `${getConfig().adminURL}` // TODO: To be deleted after local tests
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
  } catch (error) {
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
    if (showError) {
      let message
      try {
        message = errorToString(await response.json())
      } catch {
        message = 'Oops! something went wrong.'
      }
      store.dispatch(errorOccurred(message))
      throw new Error(message)
    }
  }
}
