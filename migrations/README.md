# Voltaserve Migrations

Build:

```shell
cargo build --release
```

Run locally:

```shell
DATABASE_URL=postgresql://voltaserve:voltaserve@localhost:5432/voltaserve ./target/release/migrate up
```
