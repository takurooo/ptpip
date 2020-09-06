# PTP-IP( Picture Trasfer Protocol over TCP/IP networks )

PTP-IP library for Go

# godoc
https://godoc.org/github.com/takurooo/ptpip

# goget
``
go get github.com/takurooo/ptpip
``

# Examples

```go
package main

import (
	"fmt"

	"github.com/takurooo/ptpip"
	"github.com/takurooo/ptpip/packet"
)

func main() {
    yourDeviceAddr := "192.168.3.13"

    client := ptpip.NewClient(yourDeviceAddr, nil)
    
    // establish ptp-ip connection
    if err := client.Connect(); err != nil {
        panic(err)
    }
    defer client.Disconnect()

    // GetDeviceInfo
    data, err := client.OperationRequest(0x1001, packet.DataPhaseInfoNoDataOrDataIn, 1, 0, 0, 0, 0)
    if err != nil {
        panic(err)
    }
    
    // DeviceInfo
    fmt.Println(data)

}
```