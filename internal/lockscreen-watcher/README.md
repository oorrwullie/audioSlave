# lockscreen-watcher

A lightweight macOS utility written in Swift that listens for screen lock and unlock events using Darwin notifications. Designed to be used as a subprocess in Go-based desktop automation tools.

## Purpose

This tool emits machine-readable output when the macOS screen is:

- Locked
- Unlocked

It's ideal for use cases where you need accurate detection of user presence, such as:

- Triggering smart plugs or DACs
- Controlling audio routing
- Automating security or UI behavior

## Output Format

The tool prints to stdout:
```
LOCKED
UNLOCKED
```

Each message is newline-delimited. It exits only on error or if the parent process terminates it.

## Building

This utility is managed using Swift Package Manager (SPM).

### Build the binary

```bash
cd internal/lockscreen-watcher
swift build -c release
```
The binary will be located at:
`.build/release/lockscreen-watcher`

## Integration

From your Go app, you can exec.Command() this tool and read its stdout to detect lock/unlock transitions:
```golang
cmd := exec.Command("./internal/lockscreen-watcher/.build/release/lockscreen-watcher")
stdout, _ := cmd.StdoutPipe()
// Scan output lines and act on "LOCKED"/"UNLOCKED"
```

## Requirements
	•	macOS 10.15+
	•	Swift 5.6 or higher

## License

MIT License
