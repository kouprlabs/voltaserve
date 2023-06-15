# Voltaserve Development

## Supported Operating Systems

For development, we support the following operating systems:

- Red Hat Enterprise Linux 9.x
- Rocky Linux 9.x
- AlmaLinux 9.x
- Oracle Linux 9.x

## Dependencies

Install:

```shell
curl -L "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/install.sh?t=$(date +%s)" | bash
```

Start:

```shell
curl -L "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/start.sh?t=$(date +%s)" | bash
```

Stop:

```shell
curl -L "https://raw.githubusercontent.com/kouprlabs/voltaserve/main/infra/stop.sh?t=$(date +%s)" | bash
```

_The scripts above can also be ran directly from the [infra](infra) directory._
