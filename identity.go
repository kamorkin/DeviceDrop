package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"runtime"
	"time"

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

func getDeviceID(publicKey ed25519.PublicKey) string {
	hash := sha256.Sum256(publicKey)

	return fmt.Sprintf("%x", hash)
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
		appDirectory, err := getAppDirectory()
		if err != nil {
			return nil, nil, err
		}

		deviceIDPath := filepath.Join(appDirectory, "device_id")

		_, err = os.Stat(deviceIDPath)
		if err == nil {
			fmt.Println("UUID устройства существует")
		} else if os.IsNotExist(err) {
			fmt.Println("UUID устройства отсутствует")
		} else {
			return nil, nil, err
		}

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

func createDeviceCertificate(
	publicKey ed25519.PublicKey,
	privateKey ed25519.PrivateKey,
) ([]byte, error) {
	now := time.Now()
	max := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, max)
	if err != nil {
		return nil, err
	}

	serialNumber.Add(serialNumber, big.NewInt(1))

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		NotBefore:    now,
		NotAfter:     now.AddDate(1, 0, 0),

		KeyUsage: x509.KeyUsageDigitalSignature,

		ExtKeyUsage: []x509.ExtKeyUsage{
			x509.ExtKeyUsageServerAuth,
			x509.ExtKeyUsageClientAuth,
		},
	}

	certificate, err := x509.CreateCertificate(
		rand.Reader,
		template,
		template,
		publicKey,
		privateKey,
	)
	if err != nil {
		return nil, err
	}

	return certificate, nil
}

func loadDeviceTLSCertificate() (tls.Certificate, error) {
	publicKey, privateKey, err := loadOrCreateDeviceIdentity()
	if err != nil {
		return tls.Certificate{}, err
	}

	certificateData, err := createDeviceCertificate(publicKey, privateKey)
	if err != nil {
		return tls.Certificate{}, err
	}

	certificate := tls.Certificate{
		Certificate: [][]byte{certificateData},
		PrivateKey:  privateKey,
	}

	return certificate, nil
}

func createTLSConfig() (*tls.Config, error) {
	certificate, err := loadDeviceTLSCertificate()
	if err != nil {
		return nil, err
	}

	config := &tls.Config{
		MinVersion: tls.VersionTLS13,
		Certificates: []tls.Certificate{
			certificate,
		},
	}

	return config, nil
}

func verifyDevicePublicKey(certificate *x509.Certificate,
	trustedPublicKey ed25519.PublicKey) error {
	publicKey, ok := certificate.PublicKey.(ed25519.PublicKey)
	if !ok {
		return errors.New("certificate public key is not Ed25519")
	}

	if !bytes.Equal(publicKey, trustedPublicKey) {
		return errors.New("device public key does not match trusted key")
	}

	return nil
}
