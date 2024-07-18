pub use sea_orm_migration::prelude::*;

pub mod models;

pub struct Migrator;

mod m20240718_000001_create_user;

#[async_trait::async_trait]
impl MigratorTrait for Migrator {
    fn migrations() -> Vec<Box<dyn MigrationTrait>> {
        vec![
            Box::new(m20240718_000001_create_user::Migration),
        ]
    }
}