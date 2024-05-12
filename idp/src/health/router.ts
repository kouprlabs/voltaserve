import { Router, Request, Response } from 'express'

const router = Router()

router.get('/', (_: Request, res: Response) => {
  res.send('OK')
})

export default router
