// Copyright 2024 Daniël Sonck.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.txt.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// licenses/AGPL.txt.
pub use sea_orm_migration::prelude::*;

pub mod models;

pub struct Migrator;

mod m20240718_000001_create_user;
mod m20240718_000002_create_organization;
mod m20240718_000003_create_workspace;
mod m20240718_000004_create_group;
mod m20240718_000005_create_invitation;
mod m20240718_000006_create_snapshot;
mod m20240718_000007_create_file;
mod m20240718_000008_create_task;
mod m20240718_000009_create_grouppermission;
mod m20240718_000010_create_userpermission;
mod m20240726_000001_normalize_schema;
mod m20240807_000001_add_segmentation_column;
mod m20240905_000001_add_user_active_admin_fields;

#[async_trait::async_trait]
impl MigratorTrait for Migrator {
    fn migrations() -> Vec<Box<dyn MigrationTrait>> {
        vec![
            Box::new(m20240718_000001_create_user::Migration),
            Box::new(m20240718_000002_create_organization::Migration),
            Box::new(m20240718_000003_create_workspace::Migration),
            Box::new(m20240718_000004_create_group::Migration),
            Box::new(m20240718_000005_create_invitation::Migration),
            Box::new(m20240718_000006_create_snapshot::Migration),
            Box::new(m20240718_000007_create_file::Migration),
            Box::new(m20240718_000008_create_task::Migration),
            Box::new(m20240718_000009_create_grouppermission::Migration),
            Box::new(m20240718_000010_create_userpermission::Migration),
            Box::new(m20240726_000001_normalize_schema::Migration),
            Box::new(m20240807_000001_add_segmentation_column::Migration),
            Box::new(m20240905_000001_add_user_active_admin_fields::Migration),
        ]
    }
}
