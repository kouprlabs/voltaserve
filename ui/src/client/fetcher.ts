import { getConfig } from '@/config/config'
import { getAccessToken, getAccessTokenOrRedirect } from '@/infra/token'
import store from '@/store/configure-store'
import { errorOccurred } from '@/store/ui/error'
import { errorToString } from './error'

export type FetcherOptions = {
  url: string
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH' | 'OPTIONS' | 'HEAD'
  body?: BodyInit | null
  redirect?: boolean
  authenticate?: boolean
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

export async function fetcher<T>({
  url,
  method,
  body,
  redirect,
  authenticate = true,
}: FetcherOptions) {
  const headers: HeadersInit = {
    'Content-Type': 'application/json',
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
  )
  try {
    if (response) {
      return (await response.json()) as T
    } else {
      return null
    }
  } catch {
    return null
  }
}

export async function baseFetcher(
  url: string,
  init?: RequestInit,
  redirect = true,
) {
  try {
    const response = await fetch(url, init)
    return handleFailure(response, redirect)
  } catch (error) {
    console.log('well, shit!')
  }
}

async function handleFailure(response: Response, redirect = true) {
  if (response.status <= 299) {
    return response
  } else {
    if (response.status === 401 && redirect) {
      window.location.href = '/sign-in'
    }
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
