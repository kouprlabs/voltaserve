import { ErrorCode, newError } from '@/infra/error'
import { client } from '@/infra/postgres'
import { InsertOptions, UpdateOptions, User, UserRepo } from './model'

class UserRepoImpl {
  async findByID(id: string): Promise<User> {
    const { rowCount, rows } = await client.query(
      `SELECT * FROM "user" WHERE id = $1`,
      [id],
    )
    if (rowCount < 1) {
      throw newError({
        code: ErrorCode.ResourceNotFound,
        error: `User with id=${id} not found`,
      })
    }
    return this.mapRow(rows[0])
  }

  async findByUsername(username: string): Promise<User> {
    const { rowCount, rows } = await client.query(
      `SELECT * FROM "user" WHERE username = $1`,
      [username],
    )
    if (rowCount < 1) {
      throw newError({
        code: ErrorCode.ResourceNotFound,
        error: `User with username=${username} not found`,
      })
    }
    return this.mapRow(rows[0])
  }

  async findByEmail(email: string): Promise<User> {
    const { rowCount, rows } = await client.query(
      `SELECT * FROM "user" WHERE email = $1`,
      [email],
    )
    if (rowCount < 1) {
      throw newError({
        code: ErrorCode.ResourceNotFound,
        error: `User with email=${email} not found`,
      })
    }
    return this.mapRow(rows[0])
  }

  async findByRefreshTokenValue(refreshTokenValue: string): Promise<User> {
    const { rowCount, rows } = await client.query(
      `SELECT * FROM "user" WHERE refresh_token_value = $1`,
      [refreshTokenValue],
    )
    if (rowCount < 1) {
      throw newError({
        code: ErrorCode.ResourceNotFound,
        error: `User with refresh_token_value=${refreshTokenValue} not found`,
      })
    }
    return this.mapRow(rows[0])
  }

  async findByResetPasswordToken(resetPasswordToken: string): Promise<User> {
    const { rowCount, rows } = await client.query(
      `SELECT * FROM "user" WHERE reset_password_token = $1`,
      [resetPasswordToken],
    )
    if (rowCount < 1) {
      throw newError({
        code: ErrorCode.ResourceNotFound,
        error: `User with reset_password_token=${resetPasswordToken} not found`,
      })
    }
    return this.mapRow(rows[0])
  }

  async findByEmailConfirmationToken(
    emailConfirmationToken: string,
  ): Promise<User> {
    const { rowCount, rows } = await client.query(
      `SELECT * FROM "user" WHERE email_confirmation_token = $1`,
      [emailConfirmationToken],
    )
    if (rowCount < 1) {
      throw newError({
        code: ErrorCode.ResourceNotFound,
        error: `User with email_confirmation_token=${emailConfirmationToken} not found`,
      })
    }
    return this.mapRow(rows[0])
  }

  async findByEmailUpdateToken(emailUpdateToken: string): Promise<User> {
    const { rowCount, rows } = await client.query(
      `SELECT * FROM "user" WHERE email_update_token = $1`,
      [emailUpdateToken],
    )
    if (rowCount < 1) {
      throw newError({
        code: ErrorCode.ResourceNotFound,
        error: `User with email_update_token=${emailUpdateToken} not found`,
      })
    }
    return this.mapRow(rows[0])
  }

  async findByPicture(picture: string): Promise<User> {
    const { rowCount, rows } = await client.query(
      `SELECT * FROM "user" WHERE picture = $1`,
      [picture],
    )
    if (rowCount < 1) {
      throw newError({
        code: ErrorCode.ResourceNotFound,
        error: `User with picture=${picture} not found`,
      })
    }
    return this.mapRow(rows[0])
  }

  async isUsernameAvailable(username: string): Promise<boolean> {
    const { rowCount } = await client.query(
      `SELECT * FROM "user" WHERE username = $1`,
      [username],
    )
    return rowCount === 0
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
        refresh_token_expiry,
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
        data.refreshTokenExpiry,
        data.resetPasswordToken,
        data.emailConfirmationToken,
        data.isEmailConfirmed || false,
        data.picture,
        new Date().toISOString(),
      ],
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
    const entity = await this.findByID(data.id)
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
          refresh_token_expiry = $6,
          reset_password_token = $7,
          email_confirmation_token = $8,
          is_email_confirmed = $9,
          email_update_token = $10,
          email_update_value = $11,
          picture = $12,
          update_time = $13
        WHERE id = $14
        RETURNING *`,
      [
        entity.fullName,
        entity.username,
        entity.email,
        entity.passwordHash,
        entity.refreshTokenValue,
        entity.refreshTokenExpiry,
        entity.resetPasswordToken,
        entity.emailConfirmationToken,
        entity.isEmailConfirmed,
        entity.emailUpdateToken,
        entity.emailUpdateValue,
        entity.picture,
        new Date().toISOString(),
        entity.id,
      ],
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
      refreshTokenExpiry: row.refresh_token_expiry,
      resetPasswordToken: row.reset_password_token,
      emailConfirmationToken: row.email_confirmation_token,
      isEmailConfirmed: row.is_email_confirmed,
      emailUpdateToken: row.email_update_token,
      emailUpdateValue: row.email_update_value,
      picture: row.picture,
      createTime: row.create_time,
      updateTime: row.update_time,
    }
  }
}

const userRepo: UserRepo = new UserRepoImpl()

export default userRepo
