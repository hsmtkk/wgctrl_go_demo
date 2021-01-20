package main

import(
	"fmt"
	"log"
	"golang.zx2c4.com/wireguard/wgctrl"
	"golang.zx2c4.com/wireguard/wgctrl/wgtypes"
)

func main(){
	clt,err  := wgctrl.New()	
	if err != nil {
		log.Fatal(err)
	}
	listDevices(clt)
}

func listDevices(clt *wgctrl.Client){
	devices,err := clt.Devices()
	if err != nil {
		log.Fatal(err)
	}
	for _, dev := range devices {
		printDevice(dev)
	}
}

func printDevice(dev *wgtypes.Device){
	fmt.Printf("Name: %s\n", dev.Name)
	fmt.Printf("Type: %s\n", dev.Type)
	fmt.Printf("Private key: %s\n", dev.PrivateKey)
	fmt.Printf("Public Key: %s\n", dev.PublicKey)
	fmt.Printf("Listen Port: %d\n", dev.ListenPort)
	fmt.Printf("Peers: %s\n", dev.Peers)
}
