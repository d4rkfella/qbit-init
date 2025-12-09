package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

const (
	configFilePath = "/config/qBittorrent/qBittorrent.conf"
	logFilePath    = "/config/qBittorrent/logs/qbittorrent.log"
)

const defaultConfig = `[AutoRun]
enabled=false
program=

[LegalNotice]
Accepted=true

[BitTorrent]
Session\AsyncIOThreadsCount=10
Session\DiskCacheSize=-1
Session\DiskIOReadMode=DisableOSCache
Session\DiskIOType=SimplePreadPwrite
Session\DiskIOWriteMode=EnableOSCache
Session\DiskQueueSize=4194304
Session\FilePoolSize=40
Session\HashingThreadsCount=2
Session\Port=50413
Session\ResumeDataStorageType=SQLite
Session\UseOSCache=true

[Preferences]
Connection\PortRangeMin=6881
Connection\UPnP=false
General\Locale=en
General\UseRandomPort=false
WebUI\Address=*
WebUI\AuthSubnetWhitelist=10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16
WebUI\AuthSubnetWhitelistEnabled=true
WebUI\CSRFProtection=false
WebUI\HostHeaderValidation=false
WebUI\LocalHostAuth=false
WebUI\Port=8080
WebUI\ServerDomains=*
WebUI\UseUPnP=false`

func main() {
	if err := setupConfigFile(); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up config: %v\n", err)
		os.Exit(1)
	}

	if err := setupLogSymlink(); err != nil {
		fmt.Fprintf(os.Stderr, "Error setting up log symlink: %v\n", err)
		os.Exit(1)
	}

	if err := executeQBittorrent(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing qBittorrent: %v\n", err)
		os.Exit(1)
	}
}

func setupConfigFile() error {
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Dir(configFilePath), 0755); err != nil {
			return fmt.Errorf("failed to create config directory: %w", err)
		}

		if err := os.WriteFile(configFilePath, []byte(defaultConfig), 0644); err != nil {
			return fmt.Errorf("failed to write initial config file: %w", err)
		}

		fmt.Printf("Copied default config to %s\n", configFilePath)
	}
	return nil
}

func setupLogSymlink() error {
	if err := os.MkdirAll(filepath.Dir(logFilePath), 0755); err != nil {
		return fmt.Errorf("failed to create log directory: %w", err)
	}

	if err := os.Remove(logFilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove existing log file: %w", err)
	}
	if err := os.Symlink("/proc/self/fd/1", logFilePath); err != nil {
		return fmt.Errorf("failed to create symlink to stdout: %w", err)
	}
	fmt.Printf("Symlinked log file %s to stdout\n", logFilePath)
	return nil
}

func executeQBittorrent() error {
	qbittorrentPath, err := exec.LookPath("qbittorrent-nox")
	if err != nil {
		return fmt.Errorf("failed to locate qbittorrent-nox on the system PATH: %w", err)
	}

	args := append([]string{qbittorrentPath}, os.Args[1:]...)
	err = syscall.Exec(qbittorrentPath, args, os.Environ())

	return fmt.Errorf("syscall.Exec failed to execute %s: %w", qbittorrentPath, err)
}
