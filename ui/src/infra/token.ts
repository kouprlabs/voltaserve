import TokenAPI, { Token } from '@/api/token'

export const ACCESS_TOKEN = 'voltaserve_access_token'
export const REFRESH_TOKEN = 'voltaserve_refresh_token'
export const TOKEN_EXPIRY = 'voltaserve_token_expiry'
export const COOKIE = 'voltaserve_access_token'

export async function saveToken(token: Token) {
  document.cookie = `${COOKIE}=${token.access_token}; Path=/; Max-Age=${token.expires_in}`
  localStorage.setItem(ACCESS_TOKEN, token.access_token)
  localStorage.setItem(REFRESH_TOKEN, token.refresh_token)
  const tokenExpiry = new Date()
  tokenExpiry.setSeconds(tokenExpiry.getSeconds() + token.expires_in)
  localStorage.setItem(TOKEN_EXPIRY, tokenExpiry.toISOString())
}

export async function clearToken() {
  document.cookie = `${COOKIE}=; Max-Age=-99999999;`
  localStorage.removeItem(ACCESS_TOKEN)
  localStorage.removeItem(REFRESH_TOKEN)
  localStorage.removeItem(TOKEN_EXPIRY)
}

export function getAccessTokenOrRedirect(): string {
  const accessToken = localStorage.getItem(ACCESS_TOKEN)
  const tokenExpiry = localStorage.getItem(TOKEN_EXPIRY)
  if (accessToken && tokenExpiry && new Date() < new Date(tokenExpiry)) {
    return accessToken
  } else {
    clearToken()
    window.location.href = '/sign-in'
    return ''
  }
}

setInterval(async () => {
  const refreshToken = localStorage.getItem(REFRESH_TOKEN)
  const tokenExpiry = localStorage.getItem(TOKEN_EXPIRY)
  if (tokenExpiry && refreshToken) {
    const earlyExpiry = new Date(tokenExpiry)
    earlyExpiry.setMinutes(earlyExpiry.getMinutes() - 1)
    if (new Date() >= earlyExpiry) {
      const token = await TokenAPI.exchange({
        grant_type: 'refresh_token',
        refresh_token: refreshToken,
      })
      saveToken(token)
    }
  }
}, 5000)
