// swift-tools-version: 5.10

import PackageDescription

let package = Package(
    name: "VoltaserveTests",
    platforms: [
        .iOS(.v13),
        .macOS(.v10_15)
    ],
    dependencies: [
        .package(name: "Voltaserve", url: "https://github.com/kouprlabs/voltaserve-swift.git", branch: "main"),
        .package(url: "https://github.com/Alamofire/Alamofire.git", .upToNextMajor(from: "5.9.1")),
    ],
    targets: [
        .testTarget(
            name: "VoltaserveTests",
            dependencies: ["Voltaserve"],
            path: "Sources"
        ),
    ]
)
