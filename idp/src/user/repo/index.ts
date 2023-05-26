import { DatabaseType } from '@/infra/db'
import { UserRepo } from './core'
import PostgresUserRepo from './postgres'

let userRepo: UserRepo
if (process.env.DATABASE_TYPE === DatabaseType.Postgres) {
  userRepo = new PostgresUserRepo()
} else if (process.env.DATABASE_TYPE === DatabaseType.Mongo) {
  throw new Error(`Unknown database type: ${process.env.DATABASE_TYPE}`)
}

export * from './core'

export default userRepo
