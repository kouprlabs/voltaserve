// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Request, Response, Router } from 'express'
import { client as postgres } from '@/infra/postgres.ts'
import { client as meilisearch } from '@/infra/meilisearch.ts'

const router = Router()

router.get('', async (_: Request, res: Response) => {
  if (!postgres.connected) {
    res.sendStatus(503)
    return
  }
  if (!(await meilisearch.isHealthy())) {
    res.sendStatus(503)
    return
  }
  res.send('OK')
})

export default router
