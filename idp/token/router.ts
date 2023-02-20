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
