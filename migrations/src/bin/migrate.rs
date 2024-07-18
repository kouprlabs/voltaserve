use sea_orm_migration::prelude::*;

#[async_std::main]
async fn main() {
    cli::run_cli(voltaserve_migration::Migrator).await;
}
