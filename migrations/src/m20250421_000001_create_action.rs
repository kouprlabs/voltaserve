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

use crate::models::v1::Action;

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
                    .table(Action::Table)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(Action::Id)
                            .text()
                            .primary_key(),
                    )
                    .col(
                        ColumnDef::new(Action::Prompt)
                            .text()
                            .null(),
                    )
                    .col(ColumnDef::new(Action::Files).json_binary())
                    .col(ColumnDef::new(Action::Workspaces).json_binary())
                    .col(ColumnDef::new(Action::Organizations).json_binary())
                    .col(ColumnDef::new(Action::Groups).json_binary())
                    .col(ColumnDef::new(Action::Snapshots).json_binary())
                    .col(ColumnDef::new(Action::Tasks).json_binary())
                    .col(ColumnDef::new(Action::Invitations).json_binary())
                    .col(ColumnDef::new(Action::Operations).json_binary())
                    .col(ColumnDef::new(Action::Message).text())
                    .col(
                        ColumnDef::new(Action::UserID)
                            .text()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(Action::CreateTime)
                            .text()
                            .not_null(),
                    )
                    .to_owned(),
            )
            .await?;

        manager
            .create_index(
                Index::create()
                    .name("action_create_time_idx")
                    .table(Action::Table)
                    .col(Action::CreateTime)
                    .to_owned()
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
                    .table(Action::Table)
                    .to_owned(),
            )
            .await?;

        Ok(())
    }
}
