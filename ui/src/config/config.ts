import { Config } from './types'

const config: Config = {
  apiURL: '/proxy/api/v1',
  idpURL: '/proxy/idp/v1',
}

export function getConfig(): Config {
  return config
}
