package main

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/spf13/cobra"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

// reference
// https://github.com/WireGuard/wgctrl-go/blob/d44da33e9b6bfdab1fde5aa68f662e0e65788410/doc_test.go

var ListenPort = 48574

func main() {
	cmd := &cobra.Command{
		Use:  "alpha deviceName peerIPAddress peerPort peerPublicKey",
		Args: cobra.ExactArgs(4),
		Run: func(cmd *cobra.Command, args []string) {
			deviceName := args[0]
			peerIPAddress := args[1]
			peerPort, err := strconv.Atoi(args[2])
			if err != nil {
				log.Fatal(err)
			}
			peerPublicKey := args[3]
			run(deviceName, peerIPAddress, peerPort, peerPublicKey)
		},
	}
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(deviceName string, peerIPAddress string, peerPort int, peerPublicKey string) {
	clt, err := wgctrl.New()
	if err != nil {
		log.Fatal(err)
	}
	defer clt.Close()

	peerInfo := peerInfo{
		ipAddress:        peerIPAddress,
		port:             peerPort,
		encodedPublicKey: peerPublicKey,
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
