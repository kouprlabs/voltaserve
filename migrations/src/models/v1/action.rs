use sea_orm_migration::prelude::*;

#[derive(Iden)]
pub enum Action {
    Table,
    Id,
    Prompt,
    Files,
    Workspaces,
    Organizations,
    Groups,
    Snapshots,
    Tasks,
    Invitations,
    Operations,
    Message,
    UserID,
    CreateTime,
}
