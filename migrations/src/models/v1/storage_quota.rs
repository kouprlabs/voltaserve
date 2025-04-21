use sea_orm_migration::prelude::*;

#[derive(Iden)]
pub enum StorageQuota {
    Table,
    Id,
    UserID,
    StorageCapacity,
    CreateTime,
    UpdateTime,
}
