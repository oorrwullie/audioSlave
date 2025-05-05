// swift-tools-version:5.7
import PackageDescription

let package = Package(
    name: "dac-checker",
    platforms: [.macOS(.v12)],
    targets: [
        .executableTarget(
            name: "dac-checker",
            path: "Sources/dac-checker"
        )
    ]
)
