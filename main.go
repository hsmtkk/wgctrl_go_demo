package main

import(
	"fmt"
	"log"
	"golang.zx2c4.com/wireguard/wgctrl"
)

func main(){
	clt,err  := wgctrl.New()	
	if err != nil {
		log.Fatal(err)
	}
	devices,err := clt.Devices()
	if err != nil {
		log.Fatal(err)
	}
	for _, dev := range devices {
		fmt.Printf("%v\n", dev)
	}
}