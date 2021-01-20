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
		Use:  "alpha deviceName peerIPAddress peerPort peerPublicKey allowedIPs",
		Args: cobra.ExactArgs(5),
		Run: func(cmd *cobra.Command, args []string) {
			deviceName := args[0]
			peerIPAddress := args[1]
			peerPort, err := strconv.Atoi(args[2])
			if err != nil {
				log.Fatal(err)
			}
			peerPublicKey := args[3]
			allowedIPs := args[4]
			run(deviceName, peerIPAddress, peerPort, peerPublicKey, allowedIPs)
		},
	}
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func run(deviceName string, peerIPAddress string, peerPort int, peerPublicKey string, allowedIPs string) {
	clt, err := wgctrl.New()
	if err != nil {
		log.Fatal(err)
	}
	defer clt.Close()

	peerInfo := peerInfo{
		ipAddress:        peerIPAddress,
		port:             peerPort,
		encodedPublicKey: peerPublicKey,
		allowedIPs:       allowedIPs,
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
	allowedIPs       string
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
	_, allowedIPs, err := net.ParseCIDR(peerInfo.allowedIPs)
	if err != nil {
		return fmt.Errorf("failed to parse allowed IPs; %w", err)
	}
	peer := wgtypes.PeerConfig{
		PublicKey: peerPublicKey,
		Endpoint: &net.UDPAddr{
			IP:   net.ParseIP(peerInfo.ipAddress),
			Port: peerInfo.port,
		},
		AllowedIPs: []net.IPNet{*allowedIPs},
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
