use sea_orm_migration::prelude::*;

use crate::models::v1::Snapshot;

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
                    .table(Snapshot::Table)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(Snapshot::Id)
                            .text()
                            .primary_key(),
                    )
                    .col(ColumnDef::new(Snapshot::Version).big_integer())
                    .col(ColumnDef::new(Snapshot::Original).json_binary())
                    .col(ColumnDef::new(Snapshot::Preview).json_binary())
                    .col(ColumnDef::new(Snapshot::Text).json_binary())
                    .col(ColumnDef::new(Snapshot::Ocr).json_binary())
                    .col(ColumnDef::new(Snapshot::Entities).json_binary())
                    .col(ColumnDef::new(Snapshot::Mosaic).json_binary())
                    .col(ColumnDef::new(Snapshot::Thumbnail).json_binary())
                    .col(ColumnDef::new(Snapshot::Language).text())
                    .col(ColumnDef::new(Snapshot::Status).text())
                    .col(ColumnDef::new(Snapshot::TaskId).text())
                    .col(ColumnDef::new(Snapshot::CreateTime).text())
                    .col(ColumnDef::new(Snapshot::UpdateTime).text())
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
                    .if_exists()
                    .table(Snapshot::Table)
                    .to_owned(),
            )
            .await
    }
}
