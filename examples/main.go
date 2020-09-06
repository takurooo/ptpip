package main

import (
	"fmt"

	"github.com/takurooo/ptpip"
	"github.com/takurooo/ptpip/packet"
)

func main() {
	yourDeviceAddr := "192.168.3.13"

	client := ptpip.NewClient(yourDeviceAddr, nil)
	if err := client.Connect(); err != nil {
		panic(err)
	}
	defer client.Disconnect()

	// GetDeviceInfo
	data, err := client.OperationRequest(0x1001, packet.DataPhaseInfoNoDataOrDataIn, 1, 0, 0, 0, 0, nil)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)

}
