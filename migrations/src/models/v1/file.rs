use sea_orm_migration::prelude::*;

#[derive(Iden)]
pub enum File {
    Table,
    Id,
    Name,
    Type,
    ParentId,
    WorkspaceId,
    SnapshotId,
    CreateTime,
    UpdateTime,
}

#[derive(Iden)]
pub enum SnapshotFile {
    Table,
    SnapshotId,
    FileId,
    CreateTime,
}

#[derive(Iden)]
pub enum Userpermission {
    Table,
    Id,
    UserId,
    ResourceId,
    Permission,
    CreateTime,
}

#[derive(Iden)]
pub enum Grouppermission {
    Table,
    Id,
    GroupId,
    ResourceId,
    Permission,
    CreateTime,
}
