// Copyright 2024 Daniël Sonck.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
use sea_orm_migration::prelude::*;

use crate::models::v1::{File, GroupUser, OrganizationUser, Snapshot, SnapshotFile, Workspace};

#[derive(DeriveMigrationName)]
pub struct Migration;

#[async_trait::async_trait]
impl MigrationTrait for Migration {
    async fn up(
        &self,
        manager: &SchemaManager,
    ) -> Result<(), DbErr> {
        manager
            .drop_table(
                // Remove group_user table if it still exists
                Table::drop()
                    .table(GroupUser::Table)
                    .if_exists()
                    .to_owned(),
            )
            .await?;

        manager
            .drop_table(
                // Remove organization_user if it still exists
                Table::drop()
                    .table(OrganizationUser::Table)
                    .if_exists()
                    .to_owned(),
            )
            .await?;

        manager
            .alter_table(
                Table::alter()
                    .table(Workspace::Table)
                    // Lock root folder if a workspace still actively uses it
                    .add_foreign_key(
                        TableForeignKey::new()
                            .from_tbl(Workspace::Table)
                            .from_col(Workspace::RootId)
                            .to_tbl(File::Table)
                            .to_col(File::Id),
                    )
                    .to_owned(),
            )
            .await?;

        manager
            .alter_table(
                Table::alter()
                    .table(File::Table)
                    // Orphan children if a parent is deleted
                    .add_foreign_key(
                        TableForeignKey::new()
                            .from_tbl(File::Table)
                            .from_col(File::ParentId)
                            .to_tbl(File::Table)
                            .to_col(File::Id)
                            .on_delete(ForeignKeyAction::SetNull),
                    )
                    // Create stricter file -> snapshot foreign key constraint to prevent in-use snapshots
                    // from being deleted
                    .add_foreign_key(
                        TableForeignKey::new()
                            .from_tbl(File::Table)
                            .from_col(File::SnapshotId)
                            .to_tbl(Snapshot::Table)
                            .to_col(Snapshot::Id),
                    )
                    .to_owned(),
            )
            .await?;

        // Recreate stricter snapshot_file -> snapshot foreign key constraint to prevent in-use
        // snapshots from being deleted.
        manager
            .alter_table(
                Table::alter()
                    .table(SnapshotFile::Table)
                    .drop_foreign_key(Alias::new("snapshot_file_snapshot_id_fkey"))
                    .add_foreign_key(
                        TableForeignKey::new()
                            .from_tbl(SnapshotFile::Table)
                            .from_col(SnapshotFile::SnapshotId)
                            .to_tbl(Snapshot::Table)
                            .to_col(Snapshot::Id),
                    )
                    .to_owned(),
            )
            .await?;

        Ok(())
    }
}
