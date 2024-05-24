import { Config } from './types'

const config: Config = {
  apiURL: '/proxy/api/v2',
  idpURL: '/proxy/idp/v2',
}

export function getConfig(): Config {
  return config
}
