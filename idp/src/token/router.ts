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
import { exchange, TokenExchangeOptions } from './service'

const router = Router()

router.post('/', async (req: Request, res: Response, next: NextFunction) => {
  try {
    const options = req.body as TokenExchangeOptions
    res.json(await exchange(options))
  } catch (err) {
    res.status(400)
    if (
      err.error &&
      (err.error === 'invalid_grant' ||
        err.error === 'invalid_request' ||
        err.error === 'unsupported_grant_type')
    ) {
      res.json(err)
    }
    next(err)
  }
})

export default router
