use sea_orm_migration::prelude::*;

use crate::models::v1::{
    User, Userpermission,
};

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
                    .table(Userpermission::Table)
                    .col(
                        ColumnDef::new(Userpermission::Id)
                            .text()
                            .primary_key(),
                    )
                    .col(ColumnDef::new(Userpermission::UserId).text())
                    .foreign_key(
                        ForeignKey::create()
                            .from(Userpermission::Table, Userpermission::UserId)
                            .to(User::Table, User::Id)
                            .on_delete(ForeignKeyAction::Cascade),
                    )
                    .col(ColumnDef::new(Userpermission::ResourceId).text())
                    .col(ColumnDef::new(Userpermission::Permission).text())
                    .col(
                        ColumnDef::new(Userpermission::CreateTime)
                            .text()
                            .not_null(),
                    )
                    .to_owned(),
            )
            .await?;

        manager
            .create_index(
                Index::create()
                    .table(Userpermission::Table)
                    .col(Userpermission::UserId)
                    .col(Userpermission::ResourceId)
                    .unique()
                    .to_owned(),
            )
            .await?;

        manager
            .create_index(
                Index::create()
                    .table(Userpermission::Table)
                    .col(Userpermission::UserId)
                    .to_owned(),
            )
            .await?;
        manager
            .create_index(
                Index::create()
                    .table(Userpermission::Table)
                    .col(Userpermission::ResourceId)
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
                    .table(Userpermission::Table)
                    .to_owned(),
            )
            .await?;

        Ok(())
    }
}
