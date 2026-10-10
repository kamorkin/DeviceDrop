package main

import (
	"crypto/ed25519"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/uuid"
	"github.com/zalando/go-keyring"
)

func getAppDirectory() (string, error) {
	var appDirectory string

	switch runtime.GOOS {
	case "windows":
		directoryPath := os.Getenv("LOCALAPPDATA")
		if directoryPath == "" {
			return "", fmt.Errorf("LOCALAPPDATA is not set")
		}

		appDirectory = filepath.Join(directoryPath, "DeviceDrop")

	case "darwin":
		directoryPath, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}

		appDirectory = filepath.Join(
			directoryPath,
			"Library",
			"Application Support",
			"DeviceDrop",
		)

	default:
		return "", fmt.Errorf("unsupported OS: %s", runtime.GOOS)
	}

	return appDirectory, nil
}

func loadOrCreateDeviceID() (uuid.UUID, error) {

	appDirectory, err := getAppDirectory()
	if err != nil {
		return uuid.Nil, err
	}

	err = os.MkdirAll(appDirectory, 0755)
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

func loadOrCreateDeviceIdentity() (ed25519.PublicKey, ed25519.PrivateKey, error) {
	privateKey, err := loadDevicePrivateKey()

	if errors.Is(err, keyring.ErrNotFound) {
		publicKey, privateKey, err := generateDeviceKeys()
		if err != nil {
			return nil, nil, err
		}

		err = saveDevicePrivateKey(privateKey)
		if err != nil {
			return nil, nil, err
		}

		return publicKey, privateKey, nil
	}

	if err != nil {
		return nil, nil, err
	}

	publicKey := privateKey.Public().(ed25519.PublicKey)

	return publicKey, privateKey, nil
}