// Copyright (c) 2023 Anass Bouassaba.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
use sea_orm_migration::prelude::*;

use crate::models::v1::{MurphQuota, StorageQuota};

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
                    .table(MurphQuota::Table)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(MurphQuota::Id)
                            .text()
                            .primary_key(),
                    )
                    .col(
                        ColumnDef::new(MurphQuota::UserID)
                            .text()
                            .not_null()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(MurphQuota::ActionsPerMonth)
                            .big_integer()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(MurphQuota::ActionsUsage)
                            .big_integer()
                            .not_null()
                            .default(0),
                    )
                    .col(
                        ColumnDef::new(MurphQuota::ActionsResetTime)
                            .text()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(MurphQuota::MemoryWindow)
                            .big_integer()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(MurphQuota::CreateTime)
                            .text()
                            .not_null(),
                    )
                    .col(ColumnDef::new(StorageQuota::UpdateTime).text())
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
                    .table(MurphQuota::Table)
                    .to_owned(),
            )
            .await?;

        Ok(())
    }
}
