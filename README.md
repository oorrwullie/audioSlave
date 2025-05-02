# audioSlave

**audioSlave** is a lightweight Go service that monitors your MacBook's lock/unlock state and automatically controls a HomeKit plug based on your DAC connection status.

---

## ✨ Features

- Detects Mac unlock (Wake) and lock (Sleep) events
- Checks if a specific DAC (Digital-to-Analog Converter) is connected
- Verifies that the DAC is running at the desired sample rate (e.g., 384000 Hz)
- Turns ON a HomeKit smart plug on unlock if DAC is ready
- Turns OFF the smart plug on lock
- Simple `configure` wizard to set up without editing files manually
- Runs automatically at login via LaunchAgent

---

## 🚀 Installation

### 1. Clone the repository

```bash
git clone https://github.com/oorrwullie/audioSlave.git
cd audioSlave
```

### 2. Install dependencies

Make sure you have [Go](https://golang.org/dl/) installed.

```bash
brew install go
```

### 3. Build and install

```bash
make install
```

This will:
- Build the `audioSlave` binary
- Move it to `/usr/local/bin/`
- Set up the configuration folder `/usr/local/etc/audioSlave/`
- Install and load the LaunchAgent

---

## ⚙️ First-Time Configuration

After installing, run:

```bash
audioSlave configure
```

This will guide you through:
- Selecting your DAC device
- Setting the expected sample rate
- Entering your Homebridge base URL

---

## 📋 Commands

| Command              | Description                                |
|----------------------|--------------------------------------------|
| `make install`       | Build and install the app and LaunchAgent  |
| `make uninstall`     | Remove all installed components            |
| `make clean`         | Remove local build artifacts               |
| `make configure`     | Run the configuration wizard               |
| `audioSlave`         | Start the service manually                 |

---

## 🖥 Auto Start on Boot

`audioSlave` installs a LaunchAgent:

```bash
~/Library/LaunchAgents/com.oorrwullie.audioslave.plist
```

This ensures that `audioSlave` runs every time you log into your Mac.

---

## 📜 License

MIT License

---

## ✍️ Author

Built with ❤️ by [oorrwullie](https://github.com/oorrwullie)
