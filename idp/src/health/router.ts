// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Hono } from 'hono'
import { withPostgres } from '@/infra/postgres.ts'
import { client as meilisearch } from '@/infra/meilisearch.ts'

const router = new Hono()

router.get('', async (c) => {
  try {
    await withPostgres((client) => client.queryArray('SELECT 1'))
  } catch {
    return c.body(null, 503)
  }
  if (!(await meilisearch.isHealthy())) {
    return c.body(null, 503)
  }
  return c.text('OK')
})

export default router
