use sea_orm_migration::prelude::*;

use crate::models::v1::{Group, GroupUser, Organization, User};

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
                            .not_null()
                            .default(Keyword::CurrentTimestamp),
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

        manager
            .create_table(
                Table::create()
                    .table(GroupUser::Table)
                    .col(ColumnDef::new(GroupUser::GroupId).text())
                    .foreign_key(
                        ForeignKey::create()
                            .from(GroupUser::Table, GroupUser::GroupId)
                            .to(Group::Table, Group::Id)
                            .on_delete(ForeignKeyAction::Cascade),
                    )
                    .col(ColumnDef::new(GroupUser::UserId).text())
                    .foreign_key(
                        ForeignKey::create()
                            .from(GroupUser::Table, GroupUser::UserId)
                            .to(User::Table, User::Id)
                            .on_delete(ForeignKeyAction::Cascade),
                    )
                    .col(
                        ColumnDef::new(GroupUser::CreateTime)
                            .text()
                            .not_null()
                            .default(Keyword::CurrentTimestamp),
                    )
                    .primary_key(
                        Index::create()
                            .col(GroupUser::GroupId)
                            .col(GroupUser::UserId),
                    )
                    .to_owned(),
            )
            .await?;

        manager
            .create_index(
                Index::create()
                    .table(GroupUser::Table)
                    .col(GroupUser::GroupId)
                    .to_owned(),
            )
            .await?;
        manager
            .create_index(
                Index::create()
                    .table(GroupUser::Table)
                    .col(GroupUser::UserId)
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
                    .table(GroupUser::Table)
                    .to_owned(),
            )
            .await?;

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
