package main

import (
	"fmt"
	"log"
	"net"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// reference
// https://github.com/WireGuard/wgctrl-go/blob/d44da33e9b6bfdab1fde5aa68f662e0e65788410/doc_test.go

var ListenPort = 48574

func main() {
	clt, err := wgctrl.New()
	if err != nil {
		log.Fatal(err)
	}
	defer clt.Close()

	peerInfo := peerInfo{
		ipAddress:        "192.0.2.1",
		port:             12345,
		encodedPublicKey: "C9VaGN9qYYWPi4IKnbM9uv75E6iL9pBqY+i+XjUc13o=",
	}

	configureDevice(clt, "wg0", peerInfo)

	listDevices(clt)
}

func listDevices(clt *wgctrl.Client) {
	devices, err := clt.Devices()
	if err != nil {
		log.Fatal(err)
	}
	for _, dev := range devices {
		printDevice(dev)
	}
}

func printDevice(dev *wgtypes.Device) {
	fmt.Printf("Name: %s\n", dev.Name)
	fmt.Printf("Type: %s\n", dev.Type)
	fmt.Printf("Private key: %s\n", dev.PrivateKey)
	fmt.Printf("Public Key: %s\n", dev.PublicKey)
	fmt.Printf("Listen Port: %d\n", dev.ListenPort)
	fmt.Printf("Peers: %v\n", dev.Peers)
}

type peerInfo struct {
	ipAddress        string
	port             int
	encodedPublicKey string
}

func configureDevice(clt *wgctrl.Client, deviceName string, peerInfo peerInfo) error {
	peerPublicKey, err := wgtypes.ParseKey(peerInfo.encodedPublicKey)
	if err != nil {
		return fmt.Errorf("failed to parse peer public key; %w", err)
	}
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		return fmt.Errorf("failed to generate private key; %w", err)
	}
	peer := wgtypes.PeerConfig{
		PublicKey: peerPublicKey,
		Endpoint: &net.UDPAddr{
			IP:   net.ParseIP(peerInfo.ipAddress),
			Port: peerInfo.port,
		},
	}
	config := wgtypes.Config{
		PrivateKey: &privateKey,
		ListenPort: &ListenPort,
		Peers:      []wgtypes.PeerConfig{peer},
	}
	if err := clt.ConfigureDevice(deviceName, config); err != nil {
		return fmt.Errorf("failed to configure device; %w", err)
	}
	return nil
}
