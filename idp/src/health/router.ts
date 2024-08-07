// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { Router, Request, Response, NextFunction } from 'express'
import { Client as PgClient } from 'pg'
import { getConfig } from '@/config/config'

const router = Router()

router.get('/', async (_: Request, res: Response, next: NextFunction) => {
  let pg: PgClient
  try {
    pg = new PgClient({ connectionString: getConfig().databaseURL })
    await pg.connect()
    await pg.query('SELECT 1')
    res.send('OK')
  } catch (err) {
    res.sendStatus(503)
    next(err)
  } finally {
    if (pg) {
      pg.end()
    }
  }
})

export default router
