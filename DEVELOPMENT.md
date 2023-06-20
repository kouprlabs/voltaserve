# Voltaserve Development

## Getting Started

Make sure your operating system is supported for development as described [here](#supported-operating-systems).

Run the install script and follow the instructions on the terminal:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/dev/install.sh?t=$(date +%s)" | sh -s
```

## Supported Operating Systems

For development, we support the following operating systems:

- [Red Hat Enterprise Linux 9.x](https://www.redhat.com/en/technologies/linux-platforms/enterprise-linux)
- [Rocky Linux 9.x](https://rockylinux.org)
- [AlmaLinux 9.x](https://almalinux.org)
- [Oracle Linux 9.x](https://www.oracle.com/linux)

It is recommended to have a fresh install of one of these operating systems, one way to do it is to have a freshly installed VM dedicated to Voltaserve development.

If you run your VM under [Windows WSL](https://learn.microsoft.com/en-us/windows/wsl), make sure `systemd` is enabled as described [here](https://learn.microsoft.com/en-us/windows/wsl/wsl-config#systemd-support).

We provide first class support for [Visual Studio Code](https://code.visualstudio.com) as an IDE/editor, like debugging configurations and extension recommendations for formatters, linters, etc.

## Infrastructure Services

Install:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/dev/install.sh?t=$(date +%s)" | sh -s
```

Start:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/dev/start.sh?t=$(date +%s)" | sh -s
```

Stop:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/dev/stop.sh?t=$(date +%s)" | sh -s
```

_Note: the scripts above can also be ran directly from the [infra/dev](infra/dev) directory._

## SQL Database

Create database:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/postgres/create_database.sql?t=$(date +%s)" | PGPASSWORD=postgres psql -U postgres -d postgres
```

Create tables:

```shell
curl -sSfL "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/postgres/schema.sql?t=$(date +%s)" | PGPASSWORD=postgres psql -U postgres -d voltaserve
```

_Note: the scripts above can also be ran directly from the [infra/postgres](infra/postgres) directory._

## Disable Firewall

During development, if you need to access your development enviroment remotely, you can disable the firewall so you don't need to open each port separately, this can be done as follows:

```shell
sudo systemctl disable firewalld
```

```shell
sudo systemctl stop firewalld
```

## Microservices Development

You can clone the [repository](https://github.com/kouprlabs/voltaserve) in your home directory, and run the microservices from there. One option could be to use Visual Studio Code's remote development feature as described [here](https://code.visualstudio.com/docs/remote/remote-overview) to connect to your development environment VM from your host OS.

You can read how to run each microservice in development mode here:

- [Voltaserve API](api/README.md)
- [Voltaserve UI](ui/README.md)
- [Voltaserve Identity Provider](idp/README.md)
- [Voltaserve WebDAV](webdav/README.md)
- [Voltaserve Conversion](conversion/README.md)
- [Voltaserve Language](language/README.md)
