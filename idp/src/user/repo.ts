// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
import { ErrorCode, newError } from '@/infra/error'
import { client } from '@/infra/postgres'
import { InsertOptions, UpdateOptions, User } from './model'

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

  async listAllPaginated(page: number, size: number): Promise<User[]> {
    const { rowCount, rows } = await client.query(
      `SELECT *
       FROM "user"
       ORDER BY create_time
       OFFSET $1
       LIMIT $2`,
      [(page - 1) * size, size],
    )
    if (rowCount < 1) {
      throw newError({
        code: ErrorCode.ResourceNotFound,
        error: `User list is empty`,
      })
    }
    return this.mapList(rows)
  }

  async getUserCount(): Promise<number> {
    const { rowCount, rows } = await client.query(
      `SELECT COUNT(id) as count FROM "user"`,
    )
    if (rowCount < 1) {
      throw newError({
        code: ErrorCode.ResourceNotFound,
        error: `Fatal database error (no users present in database)`,
      })
    }
    return rows[0].count
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
        is_admin,
        is_active,
        picture,
        create_time
      ) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14) RETURNING *`,
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
        data.isAdmin || false,
        data.isActive || true,
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
          is_admin = $10,
          is_active = $11,
          email_update_token = $12,
          email_update_value = $13,
          picture = $14,
          update_time = $15
        WHERE id = $16
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
        entity.isAdmin,
        entity.isActive,
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

  async suspend(id: string, suspend: boolean) :Promise<void> {
    await client.query('UPDATE "user" SET is_active = $1, refresh_token_value = null, refresh_token_expiry = null, update_time = $2 WHERE id = $3', [!suspend, new Date().toISOString(), id])
  }

  async enoughActiveAdmins() {
    const {rows } = await client.query('SELECT COUNT(*) as count FROM "user" WHERE is_admin IS TRUE AND is_active IS TRUE',  [])
    console.log('enough admins', rows[0].count)
    return rows[0].count > 1
  }

  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
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
      isAdmin: row.is_admin,
      isActive: row.is_active,
      emailUpdateToken: row.email_update_token,
      emailUpdateValue: row.email_update_value,
      picture: row.picture,
      createTime: row.create_time,
      updateTime: row.update_time,
    }
  }

  /* eslint-disable-next-line @typescript-eslint/no-explicit-any */
  private mapList(list: any): User[] {
    return list.map((user) => {
      return {
        id: user.id,
        fullName: user.full_name,
        username: user.username,
        email: user.email,
        passwordHash: user.password_hash,
        refreshTokenValue: user.refresh_token_value,
        refreshTokenExpiry: user.refresh_token_expiry,
        resetPasswordToken: user.reset_password_token,
        emailConfirmationToken: user.email_confirmation_token,
        isEmailConfirmed: user.is_email_confirmed,
        isAdmin: user.is_admin,
        isActive: user.is_active,
        emailUpdateToken: user.email_update_token,
        emailUpdateValue: user.email_update_value,
        picture: user.picture,
        createTime: user.create_time,
        updateTime: user.update_time,
      }
    })
  }
}

const userRepo: UserRepoImpl = new UserRepoImpl()

export default userRepo
