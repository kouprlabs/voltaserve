# Voltaserve Tests

## Getting Started

Prerequisites:
- Install [Swift](https://www.swift.org/) via [Xcode](https://developer.apple.com/xcode/) or [Swift Version Manager](https://github.com/kylef/swiftenv), the supported Swift version is 5.10.
- Install [SwiftFormat](https://github.com/nicklockwood/SwiftFormat).
- Install [SwiftLint](https://github.com/realm/SwiftLint).

Build Docker image:

```shell
docker build -t voltaserve/tests .
```

Run with Docker:

```shell
docker run --rm \
    -e API_HOST=host.docker.internal \
    -e IDP_HOST=host.docker.internal \
    -e USERNAME='anass@koupr.com' \
    -e PASSWORD='Passw0rd!' \
    voltaserve/tests
```

In Linux you should replace `host.docker.internal` with the host IP address, it can be found as follows:

```shell
ip route | grep default | awk '{print $3}'
```

Format code:
```
swiftformat .
```

Lint code:
```
swiftlint .
```