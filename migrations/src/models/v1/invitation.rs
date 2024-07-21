use sea_orm_migration::prelude::*;

#[derive(Iden)]
pub enum Invitation {
    Table,
    Id,
    OrganizationId,
    OwnerId,
    Email,
    Status,
    CreateTime,
    UpdateTime,
}