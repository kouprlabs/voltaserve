import { Client } from 'pg'
import { getConfig } from '@/config/config'

export const client = new Client({
  connectionString: getConfig().databaseURL,
})

client.connect()
