// Copyright (c) 2024 Daniël Sonck.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
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
    Segmentation,
    Thumbnail,
    Language,
    Summary,
    Intent,
    Status,
    TaskId,
    CreateTime,
    UpdateTime,
}
