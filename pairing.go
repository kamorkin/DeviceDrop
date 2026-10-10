package main

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/zalando/go-keyring"
)

type TrustedDevice struct {
	ID string `json:"id"`
	Name             string `json:"name"`
	KeyAlgorithm     string `json:"key_algorithm"`
	PublicKey        string `json:"public_key"`
	ClipboardSync    bool   `json:"clipboard_sync"`
	AutoReceiveFiles bool   `json:"auto_receive_files"`
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

func sendPairRequest(address string) error {
    conn, err := net.Dial("tcp", address)
    if err != nil {
        return err
    }
    defer conn.Close()

    if err := writeMessage(conn, messageTypePairRequest, nil); err != nil {
        return err
    }

    responseType, _, err := readMessage(conn)
    if err != nil {
        return err
    }

    switch responseType {
    case messageTypePairAccept:
        fmt.Println("Pairing accepted")
    case messageTypePairReject:
        fmt.Println("Pairing rejected")
    default:
        return fmt.Errorf("unexpected response: %d", responseType)
    }

    return nil
}

func generateDeviceKeys() (ed25519.PublicKey, ed25519.PrivateKey, error) {
    publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
    if err != nil {
        return nil, nil, err
    }

    return publicKey, privateKey, nil
}

func saveDevicePrivateKey(privateKey ed25519.PrivateKey) error {
	secret:= base64.StdEncoding.EncodeToString(privateKey)

	err := keyring.Set("DeviceDrop", "device-private-key", secret)
	if err != nil {
		return err
	}
	return nil
}

func loadDevicePrivateKey() (ed25519.PrivateKey, error) {
	secret, err := keyring.Get("DeviceDrop", "device-private-key")

	keyBytes, err := base64.StdEncoding.DecodeString(secret)
	if err != nil {
		return nil, err
	}
	
	if len(keyBytes) != ed25519.PrivateKeySize {
		return nil, errors.New("Invalid private key size")
	}
	return ed25519.PrivateKey(keyBytes), nil
}