// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { Router, Request, Response } from 'express'
import { exchange, TokenExchangeOptions } from './service'

const router = Router()

router.post('/', async (req: Request, res: Response) => {
  const options = req.body as TokenExchangeOptions
  res.json(await exchange(options))
})

export default router
