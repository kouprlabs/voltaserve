use sea_orm_migration::prelude::*;

#[derive(Iden)]
pub enum MurphQuota {
    Table,
    Id,
    UserID,
    ActionsPerMonth,
    ActionsUsage,
    ActionsResetTime,
    TagsPerMonth,
    TagsUsage,
    TagsResetTime,
    MemoryWindow,
    CreateTime,
    UpdateTime,
}
