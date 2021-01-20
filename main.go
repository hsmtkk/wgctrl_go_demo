package main

import (
	"encoding/base64"
	"fmt"
	"log"
	"net"
	"time"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func main() {
	clt, err := wgctrl.New()
	if err != nil {
		log.Fatal(err)
	}
	listDevices(clt)
	configureDevice(clt, "wg0")
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

func configureDevice(clt *wgctrl.Client, name string) {
	listenPort := 48574
	privaetKey, err := wgtypes.NewKey([]byte(decodeBase64("yF7wunIlxMPCeewVEGn0+oP0a5y5bgxQynF+irE4jm4=")))
	if err != nil {
		log.Fatal(err)
	}
	publicKey, err := wgtypes.NewKey([]byte(decodeBase64("aLY4suj1vczi9WRjwr8dNqnxGvaeZ0VGznacKQ4E9UI=")))
	if err != nil {
		log.Fatal(err)
	}
	endPoint := net.UDPAddr{
		IP:   net.ParseIP("192.168.11.21"),
		Port: 39814,
	}
	keepAliveInterval := 25 * time.Second
	_, subnet, err := net.ParseCIDR("10.0.0.2/32")
	if err != nil {
		log.Fatal(err)
	}
	allowedIPs := []net.IPNet{*subnet}
	peer := wgtypes.PeerConfig{
		PublicKey:                   publicKey,
		Remove:                      true,
		UpdateOnly:                  false,
		PresharedKey:                nil,
		Endpoint:                    &endPoint,
		PersistentKeepaliveInterval: &keepAliveInterval,
		ReplaceAllowedIPs:           true,
		AllowedIPs:                  allowedIPs,
	}
	cfg := wgtypes.Config{
		PrivateKey:   &privaetKey,
		ListenPort:   &listenPort,
		FirewallMark: nil,
		ReplacePeers: true,
		Peers:        []wgtypes.PeerConfig{peer},
	}
	if err := clt.ConfigureDevice(name, cfg); err != nil {
		log.Fatal(err)
	}
}

func decodeBase64(encoded string) []byte {
	bs, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		log.Fatal(err)
	}
	return bs
}
