use sea_orm_migration::prelude::*;

#[derive(Iden)]
pub enum Workspace {
    Table,
    Id,
    Name,
    OrganizationId,
    StorageCapacity,
    RootId,
    Bucket,
    CreateTime,
    UpdateTime,
}