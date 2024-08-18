// swift-tools-version: 5.10

import PackageDescription

let package = Package(
    name: "VoltaserveTests",
    platforms: [
        .iOS(.v13),
        .macOS(.v10_15)
    ],
    dependencies: [
        .package(name: "Voltaserve", url: "https://github.com/kouprlabs/voltaserve-swift.git", branch: "main")
    ],
    targets: [
        .testTarget(
            name: "VoltaserveTests",
            dependencies: ["Voltaserve"],
            path: "Sources"
        )
    ]
)
