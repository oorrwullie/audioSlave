// swift-tools-version:5.6
import PackageDescription

let package = Package(
    name: "lockscreen-watcher",
    platforms: [
        .macOS(.v10_15)
    ],
    targets: [
        .executableTarget(
            name: "lockscreen-watcher",
            dependencies: []
        )
    ]
)
