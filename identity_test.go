package main

import (
	"bytes"
	"crypto/ed25519"
	"crypto/x509"
	"testing"
)

func TestDeviceIdentityPersistence(t *testing.T) {
	publicKey1, _, err := loadOrCreateDeviceIdentity()
	if err != nil {
		t.Fatal(err)
	}

	publicKey2, _, err := loadOrCreateDeviceIdentity()
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(publicKey1, publicKey2) {
		t.Fatal("Public key changed")
	}
}

func TestDeviceCertificate(t *testing.T) {
	publicKey, privateKey, err := loadOrCreateDeviceIdentity()
	if err != nil {
		t.Fatal(err)
	}

	certificateData, err := createDeviceCertificate(publicKey, privateKey)
	if err != nil {
		t.Fatal(err)
	}

	certificate, err := x509.ParseCertificate(certificateData)
	if err != nil {
		t.Fatal(err)
	}

	certificatePublicKey, ok := certificate.PublicKey.(ed25519.PublicKey)
	if !ok {
		t.Fatal("Certificate does not contain an Ed25519 public key")
	}

	if !bytes.Equal(publicKey, certificatePublicKey) {
		t.Fatal("Certificate public key does not match device public key")
	}

	t.Log("Certificate serial number:", certificate.SerialNumber)
}
