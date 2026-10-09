package main

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/google/uuid"
)

type TrustedDevice struct {
	ID string `json:"id"`
}

type TrustedDevices struct {
	Devices []TrustedDevice `json:"devices"`
}

func encodeTrustedDevices(devices TrustedDevices) ([]byte, error) {
	data, err := json.MarshalIndent(devices, "", "  ")

	if err != nil {
		return nil, err
	}

	return data, nil
}

func decodeTrustedDevices(data []byte) (TrustedDevices, error) {
	var devices TrustedDevices

	err := json.Unmarshal(data, &devices)
	if err != nil {
		return TrustedDevices{}, err
	}

	return devices, nil
}

func loadTrustedDevices() (TrustedDevices, error) {
	appDirectory, err := getAppDirectory()
	if err != nil {
		return TrustedDevices{}, err
	}

	trustedDevicesPath := filepath.Join(appDirectory, "trusted_devices.json")

	data, err := os.ReadFile(trustedDevicesPath)
	if err != nil {
		if os.IsNotExist(err) {
			return TrustedDevices{}, nil
		}

		return TrustedDevices{}, err
	}

	return decodeTrustedDevices(data)
}

func saveTrustedDevices(devices TrustedDevices) error {
	appDirectory, err := getAppDirectory()
	if err != nil {
		return err
	}

	err = os.MkdirAll(appDirectory, 0755)
	if err != nil {
		return err
	}

	trustedDevicesPath := filepath.Join(appDirectory, "trusted_devices.json")

	data, err := encodeTrustedDevices(devices)

	if err != nil {
		return err
	}

	err = os.WriteFile(trustedDevicesPath, data, 0600)

	return err
}

func isTrustedDevice(devices TrustedDevices, deviceID string) bool {
	for _, device := range devices.Devices {
		if device.ID == deviceID {
			return true
		}
	}
	return false
}

func addTrustedDevice(deviceID string) error {
	parseID, err := uuid.Parse(deviceID)
	if err != nil {
		return err
	}

	deviceID = parseID.String()
	
	devices, err := loadTrustedDevices()
	if err != nil {
		return err
	}

	if isTrustedDevice(devices, deviceID) {
		return nil
	}

	devices.Devices = append(devices.Devices, TrustedDevice{ID: deviceID})

	return saveTrustedDevices(devices)
}
