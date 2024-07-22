use sea_orm_migration::prelude::*;

use crate::models::v1::{Organization, OrganizationUser, User};

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
                    .table(Organization::Table)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(Organization::Id)
                            .text()
                            .primary_key(),
                    )
                    .col(
                        ColumnDef::new(Organization::Name)
                            .text()
                            .not_null(),
                    )
                    .col(ColumnDef::new(Organization::CreateTime).text())
                    .col(ColumnDef::new(Organization::UpdateTime).text())
                    .to_owned(),
            )
            .await?;

        manager
            .create_table(
                Table::create()
                    .table(OrganizationUser::Table)
                    .if_not_exists()
                    .col(ColumnDef::new(OrganizationUser::OrganizationId).text())
                    .foreign_key(
                        ForeignKey::create()
                            .from(OrganizationUser::Table, OrganizationUser::OrganizationId)
                            .to(Organization::Table, Organization::Id)
                            .on_delete(ForeignKeyAction::Cascade),
                    )
                    .col(ColumnDef::new(OrganizationUser::UserId).text())
                    .foreign_key(
                        ForeignKey::create()
                            .from(OrganizationUser::Table, OrganizationUser::UserId)
                            .to(User::Table, User::Id)
                            .on_delete(ForeignKeyAction::Cascade),
                    )
                    .col(ColumnDef::new(OrganizationUser::CreateTime).text())
                    .primary_key(
                        Index::create()
                            .col(OrganizationUser::OrganizationId)
                            .col(OrganizationUser::UserId),
                    )
                    .to_owned(),
            )
            .await?;

        manager
            .create_index(
                Index::create()
                    .table(OrganizationUser::Table)
                    .col(OrganizationUser::OrganizationId)
                    .to_owned(),
            )
            .await?;
        manager
            .create_index(
                Index::create()
                    .table(OrganizationUser::Table)
                    .col(OrganizationUser::UserId)
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
                    .table(OrganizationUser::Table)
                    .to_owned(),
            )
            .await?;

        manager
            .drop_table(
                Table::drop()
                    .table(Organization::Table)
                    .to_owned(),
            )
            .await?;

        Ok(())
    }
}
