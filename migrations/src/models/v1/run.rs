use sea_orm_migration::prelude::*;

#[derive(Iden)]
pub enum Run {
    Table,
    Id,
    ActionId,
    OperationId,
    Error,
    CreateTime,
}
