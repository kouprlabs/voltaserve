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
                            .not_null(),
                    )
                    .col(ColumnDef::new(User::RefreshTokenExpiry).text())
                    .col(
                        ColumnDef::new(User::ResetPasswordToken)
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
                            .default(Keyword::CurrentTimestamp),
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