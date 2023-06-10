import { MeiliSearch } from 'meilisearch'
import { getConfig } from '@/config/config'

export const USER_SEARCH_INDEX = 'user'

const client = new MeiliSearch({ host: getConfig().search.url })
client.createIndex(USER_SEARCH_INDEX, { primaryKey: 'id' })

export default client
