use sea_orm_migration::prelude::*;

use crate::models::v1::{Task, User};

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
                    .table(Task::Table)
                    .col(
                        ColumnDef::new(Task::Id)
                            .text()
                            .primary_key(),
                    )
                    .col(
                        ColumnDef::new(Task::Name)
                            .text()
                            .not_null(),
                    )
                    .col(ColumnDef::new(Task::Error).text())
                    .col(ColumnDef::new(Task::Percentage).small_integer())
                    .col(
                        ColumnDef::new(Task::IsComplete)
                            .text()
                            .not_null()
                            .default(false),
                    )
                    .col(
                        ColumnDef::new(Task::IsIndeterminate)
                            .text()
                            .not_null()
                            .default(false),
                    )
                    .col(
                        ColumnDef::new(Task::UserId)
                            .text()
                            .not_null(),
                    )
                    .foreign_key(
                        ForeignKey::create()
                            .from(Task::Table, Task::UserId)
                            .to(User::Table, User::Id)
                            .on_delete(ForeignKeyAction::Cascade),
                    )
                    .col(ColumnDef::new(Task::Status).text())
                    .col(ColumnDef::new(Task::Payload).json_binary())
                    .col(
                        ColumnDef::new(Task::CreateTime)
                            .text()
                            .not_null(),
                    )
                    .col(ColumnDef::new(Task::UpdateTime).text())
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
                    .table(Task::Table)
                    .to_owned(),
            )
            .await?;

        Ok(())
    }
}
