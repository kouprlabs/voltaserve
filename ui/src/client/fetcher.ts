import { getConfig } from '@/config/config'
import { getAccessTokenOrRedirect } from '@/infra/token'
import store from '@/store/configure-store'
import { errorOccurred } from '@/store/ui/error'
import { errorToString } from './error'

export const apiFetcher = (options: FetcherOptions) =>
  fetcher({ ...options, url: `${getConfig().apiURL}${options.url}` })

export const idpFetcher = (options: FetcherOptions) =>
  fetcher({ ...options, url: `${getConfig().idpURL}${options.url}` })

export type FetcherOptions = {
  url: string
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH' | 'OPTIONS' | 'HEAD'
  body?: BodyInit | null
  redirect?: boolean
}

export const fetcher = ({ url, method, body, redirect }: FetcherOptions) =>
  baseFetcher(
    url,
    {
      method,
      body,
      headers: {
        'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
        'Content-Type': 'application/json',
      },
      credentials: 'include',
    },
    redirect,
  ).then(async (result) => {
    try {
      return await result.json()
      // eslint-disable-next-line
    } catch {}
  })

export const baseFetcher = async (
  url: string,
  init?: RequestInit,
  redirect = true,
) => handleFailure(await fetch(url, init), redirect)

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
