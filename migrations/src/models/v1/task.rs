use sea_orm_migration::prelude::*;

#[derive(Iden)]
pub enum Task {
    Table,
    Id,
    Name,
    Error,
    Percentage,
    IsComplete,
    IsIndeterminate,
    UserId,
    Status,
    Payload,
    CreateTime,
    UpdateTime,
}
