// Copyright 2024 Mateusz Kaźmierczak.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { Router, Request, Response, NextFunction } from 'express'

const router = Router()

router.get('/', async (_: Request, res: Response, next: NextFunction) => {
  try {
    res.json({ version: '3.0.0' })
  } catch (err) {
    res.sendStatus(503)
    next(err)
  }
})

export default router
