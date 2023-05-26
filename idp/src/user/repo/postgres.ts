import { client } from '@/infra/db'
import { ErrorCode, newError } from '@/infra/error'
import { Field, InsertOptions, UpdateOptions, User } from './core'

export default class PostgresUserRepo {
  async find(field: Field, value: any, canThrow?: boolean): Promise<User> {
    const { rowCount, rows } = await client.query(
      `SELECT * FROM "user" WHERE ${field} = $1`,
      [value]
    )
    if (rowCount < 1) {
      if (canThrow === true) {
        throw newError({
          code: ErrorCode.ResourceNotFound,
          error: `User with ${field}=${value} not found`,
        })
      } else {
        return null
      }
    }
    return this.mapRow(rows[0])
  }

  async findByPicture(picture: string): Promise<User> {
    const { rowCount, rows } = await client.query(
      `SELECT * FROM "user" WHERE picture = $1`,
      [picture]
    )
    if (rowCount < 1) {
      throw newError({
        code: ErrorCode.ResourceNotFound,
        error: `User with picture=${picture} not found`,
      })
    }
    return this.mapRow(rows[0])
  }

  async insert(data: InsertOptions): Promise<User> {
    const { rowCount, rows } = await client.query(
      `INSERT INTO "user" (
        id,
        full_name,
        username,
        email,
        password_hash,
        refresh_token_value,
        refresh_token_valid_to,
        reset_password_token,
        email_confirmation_token,
        is_email_confirmed,
        picture,
        create_time
      ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) RETURNING *`,
      [
        data.id,
        data.fullName,
        data.username,
        data.email,
        data.passwordHash,
        data.refreshTokenValue,
        data.refreshTokenValidTo,
        data.resetPasswordToken,
        data.emailConfirmationToken,
        data.isEmailConfirmed || false,
        data.picture,
        new Date().toISOString(),
      ]
    )
    if (rowCount < 1) {
      throw newError({
        code: ErrorCode.InternalServerError,
        error: `Inserting user with id=${data.id} failed`,
      })
    }
    return this.mapRow(rows[0])
  }

  async update(data: UpdateOptions): Promise<User> {
    const entity = await this.find('id', data.id)
    if (!entity) {
      throw newError({
        code: ErrorCode.InternalServerError,
        error: `User with id=${data.id} not found`,
      })
    }
    Object.assign(entity, data)
    entity.updateTime = new Date().toISOString()
    const { rowCount, rows } = await client.query(
      `UPDATE "user" 
        SET
          full_name = $1,
          username = $2,
          email = $3,
          password_hash = $4,
          refresh_token_value = $5,
          refresh_token_valid_to = $6,
          reset_password_token = $7,
          email_confirmation_token = $8,
          is_email_confirmed = $9,
          picture = $10,
          update_time = $11
        WHERE id = $12
        RETURNING *`,
      [
        entity.fullName,
        entity.username,
        entity.email,
        entity.passwordHash,
        entity.refreshTokenValue,
        entity.refreshTokenValidTo,
        entity.resetPasswordToken,
        entity.emailConfirmationToken,
        entity.isEmailConfirmed,
        entity.picture,
        new Date().toISOString(),
        entity.id,
      ]
    )
    if (rowCount < 1) {
      throw newError({
        code: ErrorCode.InternalServerError,
        error: `Inserting user with id=${data.id} failed`,
      })
    }
    return this.mapRow(rows[0])
  }

  async delete(id: string): Promise<void> {
    await client.query('DELETE FROM "user" WHERE id = $1', [id])
  }

  private mapRow(row: any): User {
    return {
      id: row.id,
      fullName: row.full_name,
      username: row.username,
      email: row.email,
      passwordHash: row.password_hash,
      refreshTokenValue: row.refresh_token_value,
      refreshTokenValidTo: row.refresh_token_valid_to,
      resetPasswordToken: row.reset_password_token,
      emailConfirmationToken: row.email_confirmation_token,
      isEmailConfirmed: row.is_email_confirmed,
      picture: row.picture,
      createTime: row.create_time,
      updateTime: row.update_time,
    }
  }
}
