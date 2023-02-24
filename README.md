<!-- markdownlint-disable MD033 MD041 -->
<p align="center">
  <img height="70" src="assets/brand.svg"/>
  <h1 align="center">Voltaserve</h1>
</p>

## Getting Started

Install [Docker](https://docs.docker.com/get-docker) and [Docker Compose](https://docs.docker.com/compose/install).

### Run for Development

```sh
docker compose up
```

You might need to wait a few minutes until all containers are up and running.

1. Navigate to [http://localhost:3000](http://localhost:3000).

2. Go to the [sign up page](http://localhost:3000/sign-up) and create an account.

3. Open MailCatcher [here](http://localhost:1080), then select the received email and click the "Confirm email" link.

4. Navigate to the [sign in page](http://localhost:3000/sign-in) and login with your credentials.

Additional instructions:

- [Voltaserve Identity Provider](idp/README.md)
- [Voltaserve API](api/README.md)
- [Voltaserve UI](ui/README.md)

### Run for Production

Add the following environment variables to the [.env](./.env) file, then change them accordingly to match your SMTP server:

```properties
SMTP_HOST=localhost
SMTP_PORT=25
SMTP_SECURE=true
SMTP_USERNAME=smtpuser
SMTP_PASSWORD=change_me
SMTP_SENDER_ADDRESS=no-reply@localhost
SMTP_SENDER_NAME='Voltaserve'
```

Update the following environment variables in [docker-compose.prod.yml](./docker-compose.prod.yml) by replacing `localhost` with your domain name:

```yaml
idp:
  environment:
    - URL=http://localhost:7000
    - WEB_URL=http://localhost:3000
api:
  environment:
    - URL=http://localhost:5000
    - WEB_URL=http://localhost:3000
```

The port `5000` is used for the web API, `7000` for the identity provider, and `3000` for the web UI. You can change them to match your preference.

Build Docker images:

```sh
docker-compose -f ./docker-compose.prod.yml build
```

Then:

```sh
docker-compose -f ./docker-compose.prod.yml up
```

Make sure all containers are up and running.

_Note: the ports `3000`, `5000` and `7000` need to be open and accessible from the outside, they can be mapped to any other ports of your choice._

_Note: here we assume that Voltaserve UI is accessible from `http://localhost:3000`, If not simply use the host and port that matches your configuration._

1. Navigate to [http://localhost:3000](http://localhost:3000). _(This depends on your configuration, see the notes above)_

2. Go to the [sign up page](http://localhost:3000/sign-up) and create an account.

3. Confirm your email.

4. Return to the [sign in page](http://localhost:3000/sign-in) and login with your credentials.

## Licensing

Voltaserve is released under the [The MIT License](./LICENSE).
