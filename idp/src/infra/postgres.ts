// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Pool, PoolClient } from 'https://deno.land/x/postgres@v0.19.3/mod.ts'
import { getConfig } from '@/config/config.ts'

const pool = new Pool(
  getConfig().database.url,
  getConfig().database.maxOpenConnections,
  true,
)

export async function withPostgres<T>(
  fn: (client: PoolClient) => Promise<T>,
): Promise<T> {
  const client = await pool.connect()
  try {
    return await fn(client)
  } finally {
    client.release()
  }
}
