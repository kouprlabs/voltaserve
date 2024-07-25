use sea_orm_migration::prelude::*;

use crate::models::v1::Organization;

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
                    .col(
                        ColumnDef::new(Organization::CreateTime)
                            .text()
                            .not_null(),
                    )
                    .col(ColumnDef::new(Organization::UpdateTime).text())
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
                    .table(Organization::Table)
                    .to_owned(),
            )
            .await?;

        Ok(())
    }
}
