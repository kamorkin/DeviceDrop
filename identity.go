package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
)

func loadOrCreateDeviceID() (uuid.UUID, error) {
	var appDirectory string

	switch runtime.GOOS {
	case "windows":
		directoryPath := os.Getenv("LOCALAPPDATA")
		if directoryPath == "" {
			return uuid.Nil, fmt.Errorf("LOCALAPPDATA is not set")
		}

		appDirectory = filepath.Join(directoryPath, "DeviceDrop")

	case "darwin":
		directoryPath, err := os.UserHomeDir()
		if err != nil {
			return uuid.Nil, err
		}

		appDirectory = filepath.Join(
			directoryPath,
			"Library",
			"Application Support",
			"DeviceDrop",
		)

	default:
		return uuid.Nil, fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	err := os.MkdirAll(appDirectory, 0755)
	if err != nil {
		return uuid.Nil, err
	}

	deviceIDPath := filepath.Join(appDirectory, "device_id")

	data, err := os.ReadFile(deviceIDPath)

	if err == nil {
		parsedID, err := uuid.Parse(string(data))
		if err != nil {
			return uuid.Nil, err
		}

		return parsedID, nil
	}

	if !os.IsNotExist(err) {
		return uuid.Nil, err
	}

	deviceID, err := uuid.NewRandom()
	if err != nil {
		return uuid.Nil, err
	}

	err = os.WriteFile(
		deviceIDPath,
		[]byte(deviceID.String()),
		0600,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return deviceID, nil
}

func getDeviceName() (string, error) {
	name, err := os.Hostname()
	if err != nil {
		return "", err
	}

	return name, nil
}
