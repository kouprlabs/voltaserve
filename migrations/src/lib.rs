pub use sea_orm_migration::prelude::*;

pub mod models;

pub struct Migrator;

mod m20240718_000001_create_user;
mod m20240718_000002_create_organization;
mod m20240718_000003_create_workspace;
mod m20240718_000004_create_group;

#[async_trait::async_trait]
impl MigratorTrait for Migrator {
    fn migrations() -> Vec<Box<dyn MigrationTrait>> {
        vec![
            Box::new(m20240718_000001_create_user::Migration),
            Box::new(m20240718_000002_create_organization::Migration),
            Box::new(m20240718_000003_create_workspace::Migration),
            Box::new(m20240718_000004_create_group::Migration),
        ]
    }
}
