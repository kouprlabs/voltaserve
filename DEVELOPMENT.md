# Voltaserve Development

## Supported Operating Systems

For development, we support the following operating systems:

- [Red Hat Enterprise Linux 9.x](https://www.redhat.com/en/technologies/linux-platforms/enterprise-linux)
- [Rocky Linux 9.x](https://rockylinux.org)
- [AlmaLinux 9.x](https://almalinux.org)
- [Oracle Linux 9.x](https://www.oracle.com/linux)

It is recommended to have a fresh install of one of these operating systems, one way to do it is to have a freshly installed VM dedicated to Voltaserve development.

If you run your VM under [Windows WSL](https://learn.microsoft.com/en-us/windows/wsl), make sure `systemd` is enabled as described [here](https://learn.microsoft.com/en-us/windows/wsl/wsl-config#systemd-support).

We provide first class support for [Visual Studio Code](https://code.visualstudio.com) as an IDE/editor, like debugging configurations and extension recommendations for formatters, linters, etc.

## Dependencies

Install:

```shell
curl -L "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/install.sh?t=$(date +%s)" | sudo bash
```

Start:

```shell
curl -L "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/start.sh?t=$(date +%s)" | sudo bash
```

Stop:

```shell
curl -L "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/stop.sh?t=$(date +%s)" | sudo bash
```

_The scripts above can also be ran directly from the [infra](infra) directory._
