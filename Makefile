# Project Variables
APP_NAME=audioSlave
GITHUB_HANDLE=oorrwullie
LAUNCH_AGENT_NAME=com.$(GITHUB_HANDLE).$(APP_NAME)
INSTALL_DIR=/usr/local/bin
CONFIG_DIR=/usr/local/etc/$(APP_NAME)
PLIST_DIR=~/Library/LaunchAgents
KEYCHAIN_SERVICE=audioSlave

# Files
PLIST_FILE=$(LAUNCH_AGENT_NAME).plist

all: build

build:
	go build -o $(APP_NAME) main.go || { echo "Build failed"; exit 1; }

install: build
	# Install binary
	sudo mv $(APP_NAME) $(INSTALL_DIR)/$(APP_NAME)
	sudo chmod +x $(INSTALL_DIR)/$(APP_NAME)

	# Install LaunchAgent
	mkdir -p $(PLIST_DIR)
	cp $(PLIST_FILE) $(PLIST_DIR)/$(PLIST_FILE)

	# Load LaunchAgent
	launchctl unload $(PLIST_DIR)/$(PLIST_FILE) || true
	launchctl load $(PLIST_DIR)/$(PLIST_FILE)

	@echo "⚙️  Now run 'make configure' to set up your configuration."

configure:
	$(INSTALL_DIR)/$(APP_NAME) configure

uninstall:
	# Unload LaunchAgent
	launchctl unload $(PLIST_DIR)/$(PLIST_FILE) || true

	# Remove installed files
	sudo rm -f $(INSTALL_DIR)/$(APP_NAME)
	rm -f $(PLIST_DIR)/$(PLIST_FILE)

	@echo "🧹 AudioSlave uninstalled (config preserved). Run 'make reset' to fully wipe data."

reset:
	@echo "⚠️  This will remove configuration and credentials. Proceeding..."
	sudo rm -rf $(CONFIG_DIR)

	@security delete-generic-password -s $(KEYCHAIN_SERVICE) 2>/dev/null || true

	@echo "🗑️  All configuration and credentials removed."

clean:
	rm -f $(APP_NAME)

.PHONY: all build install configure uninstall reset clean
