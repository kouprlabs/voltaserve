import store from '@/store/configure-store'
import { errorOccurred } from '@/store/ui/error'
import settings from '@/infra/settings'
import { getAccessTokenOrRedirect } from '@/infra/token'
import { errorToString } from './error'

export const apiFetch = async (url: string, init?: RequestInit) =>
  handleFailure(await fetch(`${settings.apiUrl}${url}`, init))

export const idpFetch = async (
  url: string,
  init?: RequestInit,
  redirect: boolean = true
) => handleFailure(await fetch(`${settings.idpUrl}${url}`, init), redirect)

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

async function handleFailure(response: Response, redirect: boolean = true) {
  if (response.status <= 299) {
    return response
  } else {
    if (response.status === 401 && redirect) {
      window.location.href = '/sign-in'
    }
    if (response.body) {
      const error = await response.json()
      store.dispatch(errorOccurred(errorToString(error)))
      throw error
    } else if (response.statusText) {
      store.dispatch(errorOccurred(response.statusText))
      throw response.statusText
    } else {
      const error = `Request failed with status ${response.status}`
      store.dispatch(errorOccurred(error))
      throw error
    }
  }
  return response
}
