use sea_orm_migration::prelude::*;

use crate::models::v1::{Invitation, Organization, User};

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
                    .table(Invitation::Table)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(Invitation::Id)
                            .text()
                            .primary_key(),
                    )
                    .col(ColumnDef::new(Invitation::OrganizationId).text())
                    .foreign_key(
                        ForeignKey::create()
                            .from(Invitation::Table, Invitation::OrganizationId)
                            .to(Organization::Table, Organization::Id)
                            .on_delete(ForeignKeyAction::Cascade),
                    )
                    .col(ColumnDef::new(Invitation::OwnerId).text())
                    .foreign_key(
                        ForeignKey::create()
                            .from(Invitation::Table, Invitation::OwnerId)
                            .to(User::Table, User::Id)
                            .on_delete(ForeignKeyAction::Cascade),
                    )
                    .col(ColumnDef::new(Invitation::Email).text())
                    .col(
                        ColumnDef::new(Invitation::Status)
                            .text()
                            .default("pending"),
                    )
                    .col(ColumnDef::new(Invitation::CreateTime).text())
                    .col(ColumnDef::new(Invitation::UpdateTime).text())
                    .to_owned(),
            )
            .await?;

        manager
            .create_index(
                Index::create()
                    .table(Invitation::Table)
                    .col(Invitation::OrganizationId)
                    .to_owned(),
            )
            .await?;
        manager
            .create_index(
                Index::create()
                    .table(Invitation::Table)
                    .col(Invitation::OwnerId)
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
                    .table(Invitation::Table)
                    .to_owned(),
            )
            .await?;

        Ok(())
    }
}
