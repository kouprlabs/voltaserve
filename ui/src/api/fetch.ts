import { getConfig } from '@/config/config'
import { getAccessTokenOrRedirect } from '@/infra/token'
import store from '@/store/configure-store'
import { errorOccurred } from '@/store/ui/error'
import { errorToString } from './error'

export const apiFetch = async (url: string, init?: RequestInit) =>
  handleFailure(await fetch(`${getConfig().apiURL}${url}`, init))

export const idpFetch = async (
  url: string,
  init?: RequestInit,
  redirect = true
) => handleFailure(await fetch(`${getConfig().idpURL}${url}`, init), redirect)

export const apiFetcher = (url: string) =>
  apiFetch(url, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
      'Content-Type': 'application/json',
    },
    credentials: 'include',
  }).then((result) => result.json())

export const idpFetcher = (url: string) =>
  idpFetch(url, {
    method: 'GET',
    headers: {
      'Authorization': `Bearer ${getAccessTokenOrRedirect()}`,
      'Content-Type': 'application/json',
    },
    credentials: 'include',
  }).then((result) => result.json())

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
