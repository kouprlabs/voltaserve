<!-- markdownlint-disable MD033 MD041 -->
<p align="center">
  <img height="70" src="assets/brand.svg"/>
  <h1 align="center">Voltaserve</h1>
</p>

## Getting Started

Install [Docker](https://docs.docker.com/get-docker) and [Docker Compose](https://docs.docker.com/compose/install).

Alternatively, you can instead install [Podman](https://podman.io) and [Podman Compose](https://github.com/containers/podman-compose#installation).

1. Build images:

```shell
docker compose build
```

2. Run:

```shell
docker compose up -d
```

Wait until the status of the following containers is `healthy`:

- `voltaserve-api`
- `voltaserve-idp`
- `voltaserve-conversion`
- `voltaserve-webdav`
- `voltaserve-language`
- `voltaserve-ui`

You can check that by running the command `docker ps` and look at the `STATUS` column.

3. Go to the **sign up page** <http://localhost:3000/sign-up> and create an account.

4. Open MailCatcher <http://localhost:11080>, select the received email and click the **confirm email** link.

5. Finally, go to the **sign in page** <http://localhost:3000/sign-in> and login with your credentials.

### Connect with WebDAV

Voltaserve supports [WebDAV](https://en.wikipedia.org/wiki/WebDAV), by default it's using the port `6000`.

To connect you can use [WinSCP](https://winscp.net) on Windows, [Cyberduck](https://cyberduck.io) on macOS and [Owlfiles](https://www.skyjos.com/owlfiles/) on mobile.

### Configuration

Update the `VOLTASERVE_HOSTNAME` environment variable in [.env](.env) file to match your hostname (it can optionally be an IP address as well):

```properties
VOLTASERVE_HOSTNAME="my-hostname"
```

Update the following environment variables in [.env](.env) file to match your SMTP server:

```properties
VOLTASERVE_SMTP_HOST="my-smtp-hostname"
VOLTASERVE_SMTP_PORT=587
VOLTASERVE_SMTP_SECURE=true
VOLTASERVE_SMTP_USERNAME="my-smtp-user"
VOLTASERVE_SMTP_PASSWORD="my-smtp-password"
VOLTASERVE_SMTP_SENDER_ADDRESS="no-reply@my-domain"
VOLTASERVE_SMTP_SENDER_NAME="Voltaserve"
```

The port `3000` is used for the UI, therefore it needs to be open and accessible from the outside. One way of doing it in Linux is by using `ufw`:

```shell
sudo ufw allow 3000
```

You can change the UI port to something else, other than `3000`, like `80` for example. This can be done by editing the `VOLTASERVE_UI_PORT` environment variable in [.env](.env) file as follows:

```properties
VOLTASERVE_UI_PORT=80
```

The port `6000` is used for WebDAV, you can change it by editing the `VOLTASERVE_WEBDAV_PORT` environment variable in [.env](.env) file as follows:

```properties
VOLTASERVE_WEBDAV_PORT=6000
```

The port needs to be open and accessible from the outside. One way of doing it in Linux is by using `ufw`:

```shell
sudo ufw allow 6000
```

Other ports can be changed as well by editing their respective environment variables in [.env](.env) file.

### Run Infrastructure Services as Multi-Node Clusters

1. Run:

```shell
docker compose -p voltaservecluster -f ./docker-compose.cluster.yml up -d
```

2. Initialize CockroachDB cluster:

```shell
docker exec -it voltaservecluster-cockroach1-1 ./cockroach init --insecure
```

3. Initialize Redis cluster:

```shell
docker run --rm -it --name=redis_cluster_init --network=voltaservecluster_default --ip=172.20.0.30 redis:7.0.8 redis-cli --cluster create 172.20.0.31:6373 172.20.0.32:6374 172.20.0.33:6375 172.20.0.34:6376 172.20.0.35:6377 172.20.0.36:6378 --cluster-replicas 1 --cluster-yes
```

4. Connect to CockroachDB with root to create a user and database:

```sql
CREATE DATABASE voltaserve;
CREATE USER voltaserve;
GRANT ALL PRIVILEGES ON DATABASE voltaserve TO voltaserve;
```

5. Connect to CockroachDB with the user `voltaserve` and database `voltaserve`, then run the following SQL script to create the database objects: [sql/schema.sql](sql/schema.sql).

## Troubleshooting

**My containers have issues starting up, what should I do?**

One reason might be that some ports are already allocated on your machine, in this case you can change the Voltaserve ports in [.env](.env) file.

**I'm not happy with `localhost`, can I change it?**

You can achieve this by changing the `VOLTASERVE_HOSTNAME` environment variable in [.env](.env) file.

It can be any IP address, like:

```properties
VOLTASERVE_HOSTNAME="192.168.1.100"
```

Or any custom hostname, like:

```properties
VOLTASERVE_HOSTNAME="my-hostname"
```

## Development

To setup a development environment for the purpose of developing and debugging Voltaserve, please read the development documentation available [here](DEVELOPMENT.md).

## Licensing

Voltaserve is released under the [The MIT License](./LICENSE.md).
