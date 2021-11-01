package main

import (
	"fmt"
	"log"
	"strings"

	"encoding/binary"
	"github.com/gofeel/gatt"
	"github.com/gofeel/gatt/linux/cmd"
)
var sp = &cmd.LESetScanParameters{
	LEScanType:           0x00,   // [0x00]: passive, 0x01: active
	LEScanInterval:       0x0010, // [0x10]: 0.625ms * 16
	LEScanWindow:         0x0010, // [0x10]: 0.625ms * 16
	OwnAddressType:       0x00,   // [0x00]: public, 0x01: random
	ScanningFilterPolicy: 0x00,   // [0x00]: accept all, 0x01: ignore non-white-listed.
}

var ClientOptions = []gatt.Option{
	gatt.LnxMaxConnections(1),
	gatt.LnxDeviceID(-1, true),
	gatt.LnxScanParameters(sp),
}

func onStateChanged(d gatt.Device, s gatt.State) {
	fmt.Println("State:", s)
	switch s {
	case gatt.StatePoweredOn:
		fmt.Println("scanning...")
		d.Scan([]gatt.UUID{}, true)
		return
	default:
		d.StopScanning()
	}
}

func onPeriphDiscovered(p gatt.Peripheral, a *gatt.Advertisement, rssi int) {
	var pid = p.ID()
	if strings.HasPrefix(pid, "A4:C1:38") {
		for _, v := range a.ServiceData {
			var temp = float64(binary.BigEndian.Uint16(v.Data[6:8])) / 10.0
			var humi = v.Data[8]
			fmt.Println(pid, " : ", temp, humi)
		}
	}
}

func main() {
	d, err := gatt.NewDevice(ClientOptions...)
	if err != nil {
		log.Fatalf("Failed to open device, err: %s\n", err)
		return
	}

	// Register handlers.
	d.Handle(gatt.PeripheralDiscovered(onPeriphDiscovered))
	d.Init(onStateChanged)
	select {}
}
