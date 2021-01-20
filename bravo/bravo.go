package main

import (
	"log"

	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func run(deviceName string, listenPort int) {
	clt, err := wgctrl.New()
	if err != nil {
		log.Fatal(err)
	}
	privateKey, err := wgtypes.GeneratePrivateKey()
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.PublicKey()
	peerConfig := wgtypes.PeerConfig{}
	config := wgtypes.Config{
		PrivateKey: &privateKey,
		ListenPort: &listenPort,
		Peers:      []wgtypes.PeerConfig{peerConfig},
	}
	if err := clt.ConfigureDevice(deviceName, config); err != nil {
		log.Fatal(err)
	}
}
