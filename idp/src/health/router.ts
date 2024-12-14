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
import { client as postgres } from '@/infra/postgres.ts'
import { client as meilisearch } from '@/infra/meilisearch.ts'

const router = new Hono()

router.get('', async (c) => {
  if (!postgres.connected) {
    c.status(503)
    return c.body(null)
  }
  if (!(await meilisearch.isHealthy())) {
    c.status(503)
    return c.body(null)
  }
  return c.text('OK')
})

export default router
