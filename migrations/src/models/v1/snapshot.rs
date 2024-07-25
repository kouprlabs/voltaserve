use sea_orm_migration::prelude::*;

#[derive(Iden)]
pub enum Snapshot {
    Table,
    Id,
    Version,
    Original,
    Preview,
    Text,
    Ocr,
    Entities,
    Mosaic,
    Thumbnail,
    Language,
    Status,
    TaskId,
    CreateTime,
    UpdateTime,
}
