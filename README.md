<!-- markdownlint-disable MD033 MD041 -->
<p align="center">
  <img height="70" src="assets/brand.svg"/>
  <h1 align="center">Voltaserve</h1>
</p>

## Getting Started

Install [Docker](https://docs.docker.com/get-docker) and [Docker Compose](https://docs.docker.com/compose/install).

### Run for Production

Update the `VOLTASERVE_HOSTNAME` environment variable in [.env](./.env) file to match your hostname (it can optionally be an IP address as well):

```properties
VOLTASERVE_HOSTNAME="my-hostname"
```

Update the following environment variables in [.env](./.env) file to match your SMTP server:

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

You can change the UI port to something else, other than `3000`, like `80` for example. This can be done by editing the `VOLTASERVE_UI_PORT` environment variable in [.env](./.env) file as follows:

```properties
VOLTASERVE_UI_PORT=80
```

Other ports can be changed as well by editing their respective environment variables in [.env](./.env) file.

Build Docker images:

```sh
docker compose build
```

Then:

```sh
docker compose up
```

Make sure all containers are up and running by checking their logs.

_Note: here you should replace `my-hostname` and `3000` with the hostname and port that matches your configuration, if you have SSL then make sure you are using `https://`._

1. Go to the **sign up page** <http://my-hostname:3000/sign-up> and create an account.

2. Confirm your email.

3. Finally, go to the **sign in page** <http://my-hostname:3000/sign-in> and login with your credentials.

### Run for Experimentation, Testing and Development

```sh
docker compose -f ./docker-compose.dev.yml up
```

Wait a few minutes until all containers are up and running. You can check that by looking at their logs.

1. Go to the **sign up page** <http://localhost:3000/sign-up> and create an account.

2. Open MailCatcher <http://localhost:11080>, select the received email and click the **confirm email** link.

3. Finally, go to the **sign in page** <http://localhost:3000/sign-in> and login with your credentials.

### Connect with WebDAV

Voltaserve supports [WebDAV](https://en.wikipedia.org/wiki/WebDAV), by default it's listening on the port `6000`.

You can change it by editing the `VOLTASERVE_WEBDAV_PORT` environment variable in [.env](./.env) file as follows:

```properties
VOLTASERVE_WEBDAV_PORT=6000
```

The port needs to be open and accessible from the outside. One way of doing it in Linux is by using `ufw`:

```shell
sudo ufw allow 6000
```

## Troubleshooting

**My containers have issues starting up, what should I do?**

One reason might be that some ports are already allocated on your machine, in this case you can change the Voltaserve ports in [.env](./.env) file.

**I'm not happy with `localhost`, can I change it?**

You can achieve this by changing the `VOLTASERVE_HOSTNAME` environment variable in [.env](./.env) file.

It can be any IP address, like:

```properties
VOLTASERVE_HOSTNAME="192.168.1.100"
```

Or any custom hostname, like:

```properties
VOLTASERVE_HOSTNAME="my-hostname"
```

**I want to start developing and contributing, how can I start?**

The UI `votaserve-ui` and Identity Providers `voltaserve-idp` are running in development/watch mode, you can simply open this repo in your favorite editor, change some code and it will reload automatically to reflect your changes.

As for the API `voltaserve-api`, once you edit the code, you should restart the container for the new changes to apply.

In case, you can also follow the instructions below for more detailed development instructions, like if you want to run the services directly on your OS without Docker. But when you do this, remember to stop the respective containers in Docker to avoid port conflicts.

Additional instructions:

- [Voltaserve API](api/README.md)
- [Voltaserve UI](ui/README.md)
- [Voltaserve Identity Provider](idp/README.md)
- [Voltaserve WebDAV](webdav/README.md)

## Licensing

Voltaserve is released under the [The MIT License](./LICENSE).
