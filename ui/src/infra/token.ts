import { Token } from '@/api/token'

export const LOCAL_STORAGE_KEY_ACESS_TOKEN = 'voltaserve_access_token'
export const COOKIE_NAME = 'voltaserve_access_token'

export async function saveAccessToken(token: Token) {
  document.cookie = `${COOKIE_NAME}=${token.access_token}; Path=/; Max-Age=${token.expires_in}`
  localStorage.setItem(LOCAL_STORAGE_KEY_ACESS_TOKEN, token.access_token)
}

export async function clearAccessToken() {
  document.cookie = `${COOKIE_NAME}=; Max-Age=-99999999;`
  localStorage.removeItem(LOCAL_STORAGE_KEY_ACESS_TOKEN)
}

export function getAccessToken() {
  return localStorage.getItem(LOCAL_STORAGE_KEY_ACESS_TOKEN)
}

export function getAccessTokenOrRedirect(): string {
  const accessToken = localStorage.getItem(LOCAL_STORAGE_KEY_ACESS_TOKEN)
  if (accessToken) {
    return accessToken
  } else {
    window.location.href = '/sign-in'
    return ''
  }
}
