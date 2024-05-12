import { getConfig } from '@/config/config'
import { Router, Request, Response, NextFunction } from 'express'
import { Client as PgClient } from 'pg'

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
    pg.end()
  }
})

export default router
