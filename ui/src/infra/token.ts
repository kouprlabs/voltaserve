export const LOCAL_STORAGE_KEY_ACESS_TOKEN = 'voltaserve_access_token'
export const COOKIE_NAME = 'voltaserve_access_token'

export async function saveAccessToken(token: string) {
  document.cookie = `${COOKIE_NAME}=${token}; path=/`
  localStorage.setItem(LOCAL_STORAGE_KEY_ACESS_TOKEN, token)
}

export async function clearAccessToken() {
  document.cookie = `${COOKIE_NAME}=; Max-Age=-99999999;`
  localStorage.removeItem(LOCAL_STORAGE_KEY_ACESS_TOKEN)
}

export function getAccessToken() {
  return typeof window !== 'undefined'
    ? localStorage.getItem(LOCAL_STORAGE_KEY_ACESS_TOKEN)
    : null
}

export function getAccessTokenOrRedirect(): string {
  if (typeof window !== 'undefined') {
    const accessToken = localStorage.getItem(LOCAL_STORAGE_KEY_ACESS_TOKEN)
    if (accessToken) {
      return accessToken
    } else {
      window.location.href = '/sign-in'
    }
  }
  return ''
}
