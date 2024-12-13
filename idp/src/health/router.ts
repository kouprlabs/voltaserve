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
import { Client as PgClient } from 'https://deno.land/x/postgres@v0.19.3/mod.ts'
import { getConfig } from '@/config/config.ts'

const router = Router()

router.get('', async (_: Request, res: Response) => {
  let pg: PgClient|undefined
  try {
    pg = new PgClient(getConfig().databaseURL)
    await pg.connect()
    await pg.queryObject('SELECT 1')
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
