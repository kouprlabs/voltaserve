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

use crate::models::v1::{Group, Organization};

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
                    .table(Group::Table)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(Group::Id)
                            .text()
                            .primary_key(),
                    )
                    .col(
                        ColumnDef::new(Group::Name)
                            .text()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(Group::OrganizationId)
                            .text()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .from(Group::Table, Group::OrganizationId)
                            .to(Organization::Table, Organization::Id),
                    )
                    .col(
                        ColumnDef::new(Group::CreateTime)
                            .text()
                            .not_null(),
                    )
                    .col(ColumnDef::new(Group::UpdateTime).text())
                    .to_owned(),
            )
            .await?;

        manager
            .create_index(
                Index::create()
                    .table(Group::Table)
                    .col(Group::OrganizationId)
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
                    .table(Group::Table)
                    .to_owned(),
            )
            .await?;

        Ok(())
    }
}
