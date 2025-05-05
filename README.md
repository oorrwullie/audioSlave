# audioSlave

**audioSlave** is a lightweight Go service that monitors your Mac’s lock/unlock state and automatically controls a HomeKit plug based on your DAC connection status.

---

## ✨ Features

- Detects macOS lock/unlock events (using `lockscreen-watcher`)
- Verifies that a specific DAC (Digital-to-Analog Converter) is:
  - Connected
  - Operating at a specific sample rate (via `dac-checker`)
- Turns **ON** a HomeKit smart plug on unlock if DAC is ready
- Turns **OFF** the plug on lock
- Easy `configure` wizard — no file editing needed
- Auto-starts on login using a LaunchAgent
- Supports log rotation with `newsyslog`

---

## 🚀 Installation

### 1. Clone the repository

```bash
git clone https://github.com/oorrwullie/audioSlave.git
cd audioSlave
```

### 2. Install dependencies

You’ll need [Go](https://golang.org/dl/) and [Swift](https://developer.apple.com/xcode/) (pre-installed on macOS).

Install Go using Homebrew if needed:

```bash
brew install go
```

### 3. Build and install

```bash
make install
```

This will:

- Build the `audioSlave`, `lockscreen-watcher`, and `dac-checker` binaries
- Install them into `/usr/local/bin`
- Copy the LaunchAgent to `~/Library/LaunchAgents`
- Install the log rotation config to `/etc/newsyslog.d/`
- Create a log directory at `/tmp/audioSlave-logs/`

---

## ⚙️ First-Time Configuration

After installing, run:

```bash
make configure
```

This will guide you through:

- Selecting your DAC device
- Setting the expected sample rate
- Providing your Homebridge base URL and credentials
- Selecting the HomeKit plug to control

Your configuration will be saved to:

```bash
~/.config/audioSlave/config.json
```

Credentials are stored securely in your macOS Keychain.

---

## 📋 Commands

| Command              | Description                                     |
|----------------------|-------------------------------------------------|
| `make install`       | Build and install the app and helpers           |
| `make configure`     | Run interactive setup wizard                    |
| `make uninstall`     | Remove binaries and LaunchAgent (preserves config) |
| `make reset`         | Delete config and credentials                   |
| `make rotate-logs`   | Manually trigger log rotation                   |
| `make clean`         | Remove build artifacts and logs                 |
| `audioSlave`         | Run the service manually (for debugging)        |

---

## 🖥 Auto Start on Login

The app installs a `LaunchAgent` here:

```bash
~/Library/LaunchAgents/com.oorrwullie.audioSlave.plist
```

To reload it manually:

```bash
launchctl bootout gui/$(id -u) ~/Library/LaunchAgents/com.oorrwullie.audioSlave.plist || true
launchctl bootstrap gui/$(id -u) ~/Library/LaunchAgents/com.oorrwullie.audioSlave.plist
```

Or simply reboot your Mac. ✅

---

## 📓 Log Files

Logs are written to:

```bash
/tmp/audioSlave-logs/audioSlave.out.log
/tmp/audioSlave-logs/audioSlave.err.log
```

To rotate logs manually:

```bash
make rotate-logs
```

---

## 📜 License

MIT License

---

## ✍️ Author

Built with ❤️ by [oorrwullie](https://github.com/oorrwullie)
