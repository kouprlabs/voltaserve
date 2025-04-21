// Copyright (c) 2024 Daniël Sonck.
//
// Use of this software is governed by the Business Source License
// included in the file LICENSE in the root of this repository.
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the GNU Affero General Public License v3.0 only, included in the file
// AGPL-3.0-only in the root of this repository.
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
mod m20240907_000001_add_user_force_change_password_field;
mod m20240913_000001_drop_segmentation_column;
mod m20241114_000001_drop_user_force_change_password_column;
mod m20241209_000001_add_user_failed_attempts_column;
mod m20241209_000001_add_user_locked_until_column;
mod m20250225_000001_add_summary_column;
mod m20250226_000001_add_intent_column;
mod m20250228_000001_drop_snapshot_status_column;
mod m20250328_000001_drop_userpermission_user_fkey;
mod m20250328_000002_drop_task_user_fkey;
mod m20250328_000003_drop_invitation_user_fkey;
mod m20250404_000001_update_password_hash_column;
mod m20250420_000001_add_user_strategy_column;
mod m20250421_000001_create_action;
mod m20250421_000002_create_run;
mod m20250421_000003_create_storage_quota;
mod m20250421_000004_create_murph_quota;

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
            Box::new(m20240907_000001_add_user_force_change_password_field::Migration),
            Box::new(m20240913_000001_drop_segmentation_column::Migration),
            Box::new(m20241114_000001_drop_user_force_change_password_column::Migration),
            Box::new(m20241209_000001_add_user_failed_attempts_column::Migration),
            Box::new(m20241209_000001_add_user_locked_until_column::Migration),
            Box::new(m20250225_000001_add_summary_column::Migration),
            Box::new(m20250226_000001_add_intent_column::Migration),
            Box::new(m20250228_000001_drop_snapshot_status_column::Migration),
            Box::new(m20250328_000001_drop_userpermission_user_fkey::Migration),
            Box::new(m20250328_000002_drop_task_user_fkey::Migration),
            Box::new(m20250328_000003_drop_invitation_user_fkey::Migration),
            Box::new(m20250404_000001_update_password_hash_column::Migration),
            Box::new(m20250420_000001_add_user_strategy_column::Migration),
            Box::new(m20250421_000001_create_action::Migration),
            Box::new(m20250421_000002_create_run::Migration),
            Box::new(m20250421_000003_create_storage_quota::Migration),
            Box::new(m20250421_000004_create_murph_quota::Migration),
        ]
    }
}
