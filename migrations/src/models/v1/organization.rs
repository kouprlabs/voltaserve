use sea_orm_migration::prelude::*;

#[derive(Iden)]
pub enum Organization {
    Table,
    Id,
    Name,
    CreateTime,
    UpdateTime,
}

#[derive(Iden)]
pub enum OrganizationUser {
    Table,
    OrganizationId,
    UserId,
    CreateTime,
}
