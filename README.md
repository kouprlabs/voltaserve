<!-- markdownlint-disable MD033 MD041 -->
<p align="center">
  <img height="70" src="assets/brand.svg"/>
  <h1 align="center">Voltaserve</h1>
</p>

<h3 align="center">Cloud Storage for Creators</h2>

Handle massive images with mosaic technology, extract insights from documents, stream large videos, protect documents with permissions and watermarks, preview 3D models. collaborate in real-time with your team.

## Getting Started

Pull images: (_recommended for most users_)

```shell
docker compose pull
```

Or, alternatively you can build the images from the source:

Build images:

```shell
docker compose build
```

Start containers:

```shell
docker compose up -d
```

Wait until the following containers are running:

- `voltaserve-api`
- `voltaserve-idp`
- `voltaserve-conversion`
- `voltaserve-webdav`
- `voltaserve-language`
- `voltaserve-mosaic`
- `voltaserve-ui`

You can check that by running the command `docker ps` and look at the `STATUS` column.

3. Go to the **sign up page** <http://localhost:3000/sign-up> and create an account.

4. Open MailCatcher <http://localhost:8025>, select the received email and click the **confirm email** link.

5. Finally, go to the **sign in page** <http://localhost:3000/sign-in> and login with your credentials.

### Connect with WebDAV

Voltaserve supports [WebDAV](https://en.wikipedia.org/wiki/WebDAV), by default it's using the port `8082`.

To connect you can use [Mountainduck](https://mountainduck.io), [Cyberduck](https://cyberduck.io), [WinSCP](https://winscp.net), [Owlfiles](https://www.skyjos.com/owlfiles) or [Rclone](https://rclone.org/webdav).

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

The port `8082` is used for WebDAV, you can change it by editing the `VOLTASERVE_WEBDAV_PORT` environment variable in [.env](.env) file as follows:

```properties
VOLTASERVE_WEBDAV_PORT=8082
```

The port needs to be open and accessible from the outside. One way of doing it in Linux is by using `ufw`:

```shell
sudo ufw allow 8082
```

Other ports can be changed as well by editing their respective environment variables in [.env](.env) file.

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

Voltaserve is released under the [GNU Affero General Public License v3.0](LICENSE.md).
