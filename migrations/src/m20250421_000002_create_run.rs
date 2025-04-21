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

use crate::models::v1::{Run, Action};

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
                    .table(Run::Table)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(Run::Id)
                            .text()
                            .primary_key(),
                    )
                    .col(
                        ColumnDef::new(Run::ActionId)
                            .text()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .from(Run::Table, Run::ActionId)
                            .to(Action::Table, Action::Id)
                            .on_delete(ForeignKeyAction::Cascade),
                    )
                    .col(
                        ColumnDef::new(Run::OperationId)
                            .text()
                            .not_null(),
                    )
                    .col(ColumnDef::new(Run::Error).text())
                    .col(
                        ColumnDef::new(Run::CreateTime)
                            .text()
                            .not_null(),
                    )
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
                    .table(Run::Table)
                    .to_owned(),
            )
            .await?;

        Ok(())
    }
}
