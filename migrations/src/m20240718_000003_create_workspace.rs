use sea_orm_migration::prelude::*;

use crate::models::v1::{Organization, Workspace};

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
                    .table(Workspace::Table)
                    .col(
                        ColumnDef::new(Workspace::Id)
                            .text()
                            .primary_key(),
                    )
                    .col(
                        ColumnDef::new(Workspace::Name)
                            .text()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(Workspace::OrganizationId)
                            .text()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .from(Workspace::Table, Workspace::OrganizationId)
                            .to(Organization::Table, Organization::Id)
                            .on_delete(ForeignKeyAction::Cascade),
                    )
                    .col(
                        ColumnDef::new(Workspace::StorageCapacity)
                            .big_integer()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(Workspace::RootId)
                            .text()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(Workspace::Bucket)
                            .text()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(Workspace::CreateTime)
                            .text()
                            .not_null(),
                    )
                    .col(ColumnDef::new(Workspace::UpdateTime).text())
                    .to_owned(),
            )
            .await?;

        manager
            .create_index(
                Index::create()
                    .name("workspace_organization_id_idx")
                    .table(Workspace::Table)
                    .col(Workspace::OrganizationId)
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
                    .table(Workspace::Table)
                    .to_owned(),
            )
            .await
    }
}
