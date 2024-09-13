// Copyright 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
use sea_orm_migration::prelude::*;

use crate::models::v1::{Snapshot};

#[derive(DeriveMigrationName)]
pub struct Migration;

#[async_trait::async_trait]
impl MigrationTrait for Migration {
    async fn up(
        &self,
        manager: &SchemaManager,
    ) -> Result<(), DbErr> {
        manager
            .alter_table(
                Table::alter()
                    .table(Snapshot::Table)
                    .drop_column(Snapshot::Segmentation)
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
            .alter_table(
                Table::alter()
                    .table(Snapshot::Table)
                    .add_column(ColumnDef::new(Snapshot::Segmentation).json_binary())
                    .to_owned(),
            )
            .await?;

        Ok(())
    }
}
