import { MeiliSearch } from 'meilisearch'
import { getConfig } from '@/config/config'

const client = new MeiliSearch({ host: getConfig().search.url })

export const USER_SEARCH_INDEX = 'user'

export default client
