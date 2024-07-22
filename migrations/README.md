# Voltaserve Migrations

Build:

```shell
cargo build --release
```

Run locally:

```shell
DATABASE_URL=postgresql://voltaserve@localhost:26257/voltaserve ./target/release/migrate
```