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

use crate::models::v1::{Organization, Workspace};

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
                    .table(Workspace::Table)
                    .if_not_exists()
                    .col(
                        ColumnDef::new(Workspace::Id)
                            .text()
                            .primary_key(),
                    )
                    .col(
                        ColumnDef::new(Workspace::Name)
                            .text()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(Workspace::OrganizationId)
                            .text()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .from(Workspace::Table, Workspace::OrganizationId)
                            .to(Organization::Table, Organization::Id)
                            .on_delete(ForeignKeyAction::Cascade),
                    )
                    .col(
                        ColumnDef::new(Workspace::StorageCapacity)
                            .big_integer()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(Workspace::RootId)
                            .text()
                            .unique_key(),
                    )
                    .col(
                        ColumnDef::new(Workspace::Bucket)
                            .text()
                            .not_null(),
                    )
                    .col(
                        ColumnDef::new(Workspace::CreateTime)
                            .text()
                            .not_null(),
                    )
                    .col(ColumnDef::new(Workspace::UpdateTime).text())
                    .to_owned(),
            )
            .await?;

        manager
            .create_index(
                Index::create()
                    .name("workspace_organization_id_idx")
                    .if_not_exists()
                    .table(Workspace::Table)
                    .col(Workspace::OrganizationId)
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
                    .table(Workspace::Table)
                    .to_owned(),
            )
            .await
    }
}
