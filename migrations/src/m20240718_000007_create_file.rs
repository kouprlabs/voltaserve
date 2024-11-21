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
                    .col(ColumnDef::new(File::WorkspaceId).text())
                    .foreign_key(
                        ForeignKey::create()
                            .from(File::Table, File::WorkspaceId)
                            .to(Workspace::Table, Workspace::Id)
                            .on_delete(ForeignKeyAction::Cascade),
                    )
                    .col(ColumnDef::new(File::SnapshotId).text())
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
                    .name("file_parent_id_idx")
                    .if_not_exists()
                    .table(File::Table)
                    .col(File::ParentId)
                    .to_owned(),
            )
            .await?;
        manager
            .create_index(
                Index::create()
                    .name("file_workspace_id_idx")
                    .if_not_exists()
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
                    .if_not_exists()
                    .col(ColumnDef::new(SnapshotFile::SnapshotId).text())
                    .foreign_key(
                        ForeignKey::create()
                            .from(SnapshotFile::Table, SnapshotFile::SnapshotId)
                            .to(Snapshot::Table, Snapshot::Id)
                            .on_delete(ForeignKeyAction::Cascade),
                    )
                    .col(ColumnDef::new(SnapshotFile::FileId).text())
                    .foreign_key(
                        ForeignKey::create()
                            .from(SnapshotFile::Table, SnapshotFile::FileId)
                            .to(File::Table, File::Id)
                            .on_delete(ForeignKeyAction::Cascade),
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
