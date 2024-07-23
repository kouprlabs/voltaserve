use sea_orm_migration::prelude::*;

use crate::models::v1::{File, Snapshot, SnapshotFile, Workspace};

#[derive(DeriveMigrationName)]
pub struct Migration;

#[async_trait::async_trait]
impl MigrationTrait for Migration {
    async fn up(
        &self,
        manager: &SchemaManager,
    ) -> Result<(), DbErr> {
        Self::create_file_table(manager).await?;

        Self::create_file_snapshot_table(manager).await?;

        Ok(())
    }

    async fn down(
        &self,
        manager: &SchemaManager,
    ) -> Result<(), DbErr> {
        manager
            .drop_table(
                Table::drop()
                    .table(SnapshotFile::Table)
                    .to_owned(),
            )
            .await?;

        manager
            .drop_table(
                Table::drop()
                    .table(File::Table)
                    .to_owned(),
            )
            .await?;

        Ok(())
    }
}

impl Migration {
    async fn create_file_table(manager: &SchemaManager<'_>) -> Result<(), DbErr> {
        manager
            .create_table(
                Table::create()
                    .table(File::Table)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(File::Id)
                            .text()
                            .primary_key(),
                    )
                    .col(
                        ColumnDef::new(File::Name)
                            .text()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(File::Type)
                            .text()
                            .not_null(),
                    )
                    .col(ColumnDef::new(File::ParentId).text())
                    .foreign_key(
                        ForeignKey::create()
                            .from(File::Table, File::ParentId)
                            .to(File::Table, File::Id)
                            .on_delete(ForeignKeyAction::SetNull),
                    )
                    .col(ColumnDef::new(File::WorkspaceId).text())
                    .foreign_key(
                        ForeignKey::create()
                            .from(File::Table, File::WorkspaceId)
                            .to(Workspace::Table, Workspace::Id)
                            .on_delete(ForeignKeyAction::Cascade),
                    )
                    .col(ColumnDef::new(File::SnapshotId).text())
                    .foreign_key(
                        ForeignKey::create()
                            .from(File::Table, File::SnapshotId)
                            .to(Snapshot::Table, Snapshot::Id),
                    )
                    .col(
                        ColumnDef::new(File::CreateTime)
                            .text()
                            .not_null(),
                    )
                    .col(ColumnDef::new(File::UpdateTime).text())
                    .to_owned(),
            )
            .await?;

        manager
            .create_index(
                Index::create()
                    .table(File::Table)
                    .col(File::ParentId)
                    .to_owned(),
            )
            .await?;
        manager
            .create_index(
                Index::create()
                    .table(File::Table)
                    .col(File::WorkspaceId)
                    .to_owned(),
            )
            .await?;
        Ok(())
    }

    async fn create_file_snapshot_table(manager: &SchemaManager<'_>) -> Result<(), DbErr> {
        manager
            .create_table(
                Table::create()
                    .table(SnapshotFile::Table)
                    .col(ColumnDef::new(SnapshotFile::SnapshotId).text())
                    .foreign_key(
                        ForeignKey::create()
                            .from(SnapshotFile::Table, SnapshotFile::SnapshotId)
                            .to(Snapshot::Table, Snapshot::Id),
                    )
                    .col(ColumnDef::new(SnapshotFile::FileId).text())
                    .foreign_key(
                        ForeignKey::create()
                            .from(SnapshotFile::Table, SnapshotFile::FileId)
                            .to(File::Table, File::Id),
                    )
                    .col(
                        ColumnDef::new(SnapshotFile::CreateTime)
                            .text()
                            .not_null(),
                    )
                    .primary_key(
                        Index::create()
                            .col(SnapshotFile::SnapshotId)
                            .col(SnapshotFile::FileId),
                    )
                    .to_owned(),
            )
            .await?;

        manager
            .create_index(
                Index::create()
                    .table(SnapshotFile::Table)
                    .col(SnapshotFile::SnapshotId)
                    .to_owned(),
            )
            .await?;
        manager
            .create_index(
                Index::create()
                    .table(SnapshotFile::Table)
                    .col(SnapshotFile::FileId)
                    .to_owned(),
            )
            .await?;
        Ok(())
    }
}
