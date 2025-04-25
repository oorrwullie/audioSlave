# Project Variables
APP_NAME=audioSlave
GITHUB_HANDLE=oorrwullie
LAUNCH_AGENT_NAME=com.$(GITHUB_HANDLE).$(APP_NAME)
INSTALL_DIR=/usr/local/bin
CONFIG_DIR=/usr/local/etc/$(APP_NAME)
PLIST_DIR=~/Library/LaunchAgents

# Files
CONFIG_FILE=config.json
PLIST_FILE=$(LAUNCH_AGENT_NAME).plist

# Targets
all: build

build:
\tgo build -o $(APP_NAME) main.go

install: build
\t# Install binary
\tsudo mv $(APP_NAME) $(INSTALL_DIR)/$(APP_NAME)
\tsudo chmod +x $(INSTALL_DIR)/$(APP_NAME)

\t# Install config
\tsudo mkdir -p $(CONFIG_DIR)
\tsudo cp $(CONFIG_FILE) $(CONFIG_DIR)/$(CONFIG_FILE)

\t# Install LaunchAgent
\tmkdir -p $(PLIST_DIR)
\tcp $(PLIST_FILE) $(PLIST_DIR)/$(PLIST_FILE)

\t# Load LaunchAgent
\tlaunchctl unload $(PLIST_DIR)/$(PLIST_FILE) || true
\tlaunchctl load $(PLIST_DIR)/$(PLIST_FILE)

uninstall:
\t# Unload LaunchAgent
\tlaunchctl unload $(PLIST_DIR)/$(PLIST_FILE) || true

\t# Remove installed files
\tsudo rm -f $(INSTALL_DIR)/$(APP_NAME)
\tsudo rm -rf $(CONFIG_DIR)
\trm -f $(PLIST_DIR)/$(PLIST_FILE)

clean:
\trm -f $(APP_NAME)

.PHONY: all build install uninstall clean
