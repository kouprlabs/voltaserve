use sea_orm_migration::prelude::*;

#[derive(Iden)]
pub enum Group {
    Table,
    Id,
    Name,
    OrganizationId,
    CreateTime,
    UpdateTime,
}
