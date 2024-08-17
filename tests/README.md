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
docker run --rm voltaserve/tests
```

Format code:
```
swiftformat .
```

Lint code:
```
swiftlint .
```