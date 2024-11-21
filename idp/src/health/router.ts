// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
import { Router, Request, Response } from 'express'
import { Client as PgClient } from 'pg'
import { getConfig } from '@/config/config'

const router = Router()

router.get('', async (_: Request, res: Response) => {
  let pg: PgClient
  try {
    pg = new PgClient({ connectionString: getConfig().databaseURL })
    await pg.connect()
    await pg.query('SELECT 1')
    res.send('OK')
  } catch {
    res.sendStatus(503)
  } finally {
    if (pg) {
      await pg.end()
    }
  }
})

export default router
