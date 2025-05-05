# Project Variables
APP_NAME=audioSlave
GITHUB_HANDLE=oorrwullie
LAUNCH_AGENT_NAME=com.$(GITHUB_HANDLE).$(APP_NAME)
INSTALL_DIR=/usr/local/bin
CONFIG_DIR=$(HOME)/.config/$(APP_NAME)
PLIST_DIR=$(HOME)/Library/LaunchAgents
LOG_DIR=/tmp/$(APP_NAME)-logs
KEYCHAIN_SERVICE=audioSlave

# Files
PLIST_FILE=$(LAUNCH_AGENT_NAME).plist
LOCKSCREEN_WATCHER_DIR=internal/lockscreen-watcher
LOCKSCREEN_WATCHER_BIN=$(LOCKSCREEN_WATCHER_DIR)/.build/release/lockscreen-watcher
DAC_CHECKER_DIR=internal/dac-checker
DAC_CHECKER_BIN=$(DAC_CHECKER_DIR)/.build/release/dac-checker

all: build

build: build-go build-lockscreen-watcher build-dac-checker

build-go:
	go build -o $(APP_NAME) main.go || { echo "Go build failed"; exit 1; }

build-lockscreen-watcher:
	@command -v swift >/dev/null 2>&1 || { echo >&2 "Swift is not installed. Aborting."; exit 1; }
	cd $(LOCKSCREEN_WATCHER_DIR) && swift build -c release || { echo "Swift build failed"; exit 1; }

build-dac-checker:
	@command -v swift >/dev/null 2>&1 || { echo >&2 "Swift is not installed. Aborting."; exit 1; }
	cd $(DAC_CHECKER_DIR) && swift build -c release || { echo "Swift build failed"; exit 1; }

install: build
	# Install main binary
	sudo mv $(APP_NAME) $(INSTALL_DIR)/$(APP_NAME)
	sudo chmod +x $(INSTALL_DIR)/$(APP_NAME)

	# Install lockscreen-watcher binary
	sudo cp $(LOCKSCREEN_WATCHER_BIN) $(INSTALL_DIR)/lockscreen-watcher
	sudo chmod +x $(INSTALL_DIR)/lockscreen-watcher

	# Install dac-checker binary
	sudo cp $(DAC_CHECKER_BIN) $(INSTALL_DIR)/dac-checker
	sudo chmod +x $(INSTALL_DIR)/dac-checker

	# Install LaunchAgent
	mkdir -p $(PLIST_DIR)
	cp $(PLIST_FILE) $(PLIST_DIR)/$(PLIST_FILE)

	# Ensure log directory exists
	sudo mkdir -p $(LOG_DIR)
	sudo touch $(LOG_DIR)/audioSlave.out.log $(LOG_DIR)/audioSlave.err.log
	sudo chown $$(id -u):$$(id -g) $(LOG_DIR)/*.log

	# Install log rotation config
	sudo cp contrib/audioSlave.conf /etc/newsyslog.d/audioSlave.conf

	# Load LaunchAgent (modern method)
	- launchctl bootout gui/$(shell id -u) $(PLIST_DIR)/$(PLIST_FILE) || true
	launchctl bootstrap gui/$(shell id -u) $(PLIST_DIR)/$(PLIST_FILE)

	@echo "✅ Installed AudioSlave to $(INSTALL_DIR)/$(APP_NAME)"
	@echo "✅ Installed lockscreen-watcher to $(INSTALL_DIR)/lockscreen-watcher"
	@echo "✅ Installed dac-checker to $(INSTALL_DIR)/dac-checker"
	@echo "⚙️  Now run 'make configure' to set up your configuration."

configure:
	$(INSTALL_DIR)/$(APP_NAME) configure || true
	@echo ""
	@echo "📌 To run AudioSlave automatically on login:"
	@echo "   launchctl bootstrap gui/$$(id -u) $(PLIST_DIR)/$(PLIST_FILE)"
	@echo ""
	@echo "Or simply reboot your Mac — it will start automatically."

rotate-logs:
	sudo newsyslog -F
	@echo "🔁 Forced log rotation complete. See rotated logs in $(LOG_DIR)"

uninstall:
	# Unload LaunchAgent
	- launchctl bootout gui/$(shell id -u) $(PLIST_DIR)/$(PLIST_FILE) || true

	# Remove installed binaries
	sudo rm -f $(INSTALL_DIR)/$(APP_NAME)
	sudo rm -f $(INSTALL_DIR)/lockscreen-watcher
	sudo rm -f $(INSTALL_DIR)/dac-checker
	rm -f $(PLIST_DIR)/$(PLIST_FILE)

	# Remove log config
	sudo rm -f /etc/newsyslog.d/audioSlave.conf

	@echo "🧹 AudioSlave, lockscreen-watcher, and dac-checker uninstalled (config preserved). Run 'make reset' to fully wipe data."

reset:
	@echo "⚠️  This will remove configuration and credentials. Proceeding..."
	rm -rf $(CONFIG_DIR)
	@security delete-generic-password -s $(KEYCHAIN_SERVICE) 2>/dev/null || true
	@echo "🗑️  All configuration and credentials removed."

clean:
	rm -f $(APP_NAME)
	cd $(LOCKSCREEN_WATCHER_DIR) && swift package clean
	cd $(DAC_CHECKER_DIR) && swift package clean
	rm -rf $(LOG_DIR)

.PHONY: all build build-go build-lockscreen-watcher build-dac-checker install configure uninstall reset clean rotate-logs
