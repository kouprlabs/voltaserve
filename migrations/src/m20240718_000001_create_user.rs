// Copyright (c) 2024 Daniël Sonck.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
use sea_orm_migration::prelude::*;

use crate::models::v1::User;

#[derive(DeriveMigrationName)]
pub struct Migration;

#[async_trait::async_trait]
impl MigrationTrait for Migration {
    async fn up(
        &self,
        manager: &SchemaManager,
    ) -> Result<(), DbErr> {
        manager
            .create_table(
                Table::create()
                    .table(User::Table)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(User::Id)
                            .text()
                            .primary_key(),
                    )
                    .col(
                        ColumnDef::new(User::FullName)
                            .text()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(User::Username)
                            .text()
                            .not_null()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(User::Email)
                            .text()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(User::PasswordHash)
                            .text()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(User::RefreshTokenValue)
                            .text()
                            .unique_key(),
                    )
                    .col(ColumnDef::new(User::RefreshTokenExpiry).text())
                    .col(
                        ColumnDef::new(User::ResetPasswordToken)
                            .text()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(User::EmailConfirmationToken)
                            .text()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(User::EmailUpdateToken)
                            .text()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(User::EmailUpdateValue)
                            .text()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(User::IsEmailConfirmed)
                            .boolean()
                            .not_null()
                            .default(false),
                    )
                    .col(ColumnDef::new(User::Picture).text())
                    .col(
                        ColumnDef::new(User::CreateTime)
                            .text()
                            .not_null(),
                    )
                    .col(ColumnDef::new(User::UpdateTime).text())
                    .to_owned(),
            )
            .await?;

        Ok(())
    }

    async fn down(
        &self,
        manager: &SchemaManager,
    ) -> Result<(), DbErr> {
        manager
            .drop_table(
                Table::drop()
                    .table(User::Table)
                    .to_owned(),
            )
            .await
    }
}
