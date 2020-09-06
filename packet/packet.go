package packet

import (
	"fmt"
	"io"
)

type PTPIPConn interface {
	io.Reader
	io.Writer
}

const (
	DataPhaseInfoUnkownData     uint32 = 0x00000000
	DataPhaseInfoNoDataOrDataIn uint32 = 0x00000000
	DataPhaseInfoDataOut        uint32 = 0x00000000
)

// Packet Type
const (
	PacketTypeInitCommandRequest uint32 = 0x00000001
	PacketTypeInitCommandAck     uint32 = 0x00000002
	PacketTypeInitEventRequest   uint32 = 0x00000003
	PacketTypeInitEventAck       uint32 = 0x00000004
	PacketTypeInitFail           uint32 = 0x00000005
	PacketTypeOperationRequest   uint32 = 0x00000006
	PacketTypeOperationResponse  uint32 = 0x00000007
	PacketTypeEvent              uint32 = 0x00000008
	PacketTypeStartData          uint32 = 0x00000009
	PacketTypeData               uint32 = 0x0000000A
	PacketTypeCancel             uint32 = 0x0000000B
	PacketTypeEndData            uint32 = 0x0000000C
	PacketTypeProbeRequest       uint32 = 0x0000000D
	PacketTypeProbeResponse      uint32 = 0x0000000E
)

const (
	ResponseCodeUndefined                             uint16 = 0x2000
	ResponseCodeOK                                    uint16 = 0x2001
	ResponseCodeGeneralError                          uint16 = 0x2002
	ResponseCodeSessionNotOpen                        uint16 = 0x2003
	ResponseCodeInvalidTransactionID                  uint16 = 0x2004
	ResponseCodeOperationNotSupported                 uint16 = 0x2005
	ResponseCodePrameterNotSupported                  uint16 = 0x2006
	ResponseCodeIncompleteTransfer                    uint16 = 0x2007
	ResponseCodeInvalidStorageID                      uint16 = 0x2008
	ResponseCodeInvalidObjectHandle                   uint16 = 0x2009
	ResponseCodeDevicePropNotSupported                uint16 = 0x200A
	ResponseCodeInvalidObjectFormatCode               uint16 = 0x200B
	ResponseCodeStoreFull                             uint16 = 0x200C
	ResponseCodeObjectWriteProtected                  uint16 = 0x200D
	ResponseCodeStoreReadOnly                         uint16 = 0x200E
	ResponseCodeAccessDenied                          uint16 = 0x200F
	ResponseCodeNoThumbnailPresent                    uint16 = 0x2010
	ResponseCodeSelfTestFailed                        uint16 = 0x2011
	ResponseCodePartialDelection                      uint16 = 0x2012
	ResponseCodeStoreNotAvailable                     uint16 = 0x2013
	ResponseCodeSpecificationByFormatUnsupported      uint16 = 0x2014
	ResponseCodeNoValidObjectInfo                     uint16 = 0x2015
	ResponseCodeInvalidCodeFormat                     uint16 = 0x2016
	ResponseCodeUnknownVendorCode                     uint16 = 0x2017
	ResponseCodeCaptureAlreadyTerminated              uint16 = 0x2018
	ResponseCodeDeviceBusy                            uint16 = 0x2019
	ResponseCodeInvalidParentObject                   uint16 = 0x201A
	ResponseCodeInvalidDevicePropFormat               uint16 = 0x201B
	ResponseCodeInvalidDevicePropValue                uint16 = 0x201C
	ResponseCodeInvalidParameter                      uint16 = 0x201D
	ResponseCodeSessionAlreadyOpen                    uint16 = 0x201E
	ResponseCodeTransactionCancelled                  uint16 = 0x201F
	ResponseCodeSpecificationOfDestinationUnsupported uint16 = 0x201F
)

// InitCommandRequestPacket ...
type InitCommandRequestPacket struct {
	GUID            []byte
	FriendlyName    string
	ProtocolVersion uint32
}

func (cr InitCommandRequestPacket) String() string {
	var s string
	s += fmt.Sprintf("----------------\n")
	s += fmt.Sprintf("InitCommandRequestPacket\n")
	s += fmt.Sprintf("----------------\n")
	s += fmt.Sprintf("GUID             : %v\n", cr.GUID)
	s += fmt.Sprintf("FriendlyName     : %v\n", cr.FriendlyName)
	s += fmt.Sprintf("ProtocolVersion  : %v", cr.ProtocolVersion)
	return s
}

// InitCommandAckPacket ...
type InitCommandAckPacket struct {
	ConnectionNumber uint32
	GUID             []byte
	FriendlyName     string
	ProtocolVersion  uint32
}

func (ca InitCommandAckPacket) String() string {
	var s string
	s += fmt.Sprintf("----------------\n")
	s += fmt.Sprintf("InitCommandAckPacket\n")
	s += fmt.Sprintf("----------------\n")
	s += fmt.Sprintf("ConnectionNumber : 0x%08x\n", ca.ConnectionNumber)
	s += fmt.Sprintf("GUID             : %v\n", ca.GUID)
	s += fmt.Sprintf("FriendlyName     : %v\n", ca.FriendlyName)
	s += fmt.Sprintf("ProtocolVersion  : %v", ca.ProtocolVersion)
	return s
}

// OperationRequestPacket ...
type OperationRequestPacket struct {
	DataPhaseInfo uint32
	OperationCode uint16
	TransactionID uint32
	P1            uint32
	P2            uint32
	P3            uint32
	P4            uint32
}

func (o OperationRequestPacket) String() string {
	var s string
	s += fmt.Sprintf("----------------\n")
	s += fmt.Sprintf("OperationRequestPacket\n")
	s += fmt.Sprintf("----------------\n")
	s += fmt.Sprintf("DataPhaseInfo    : 0x%08x\n", o.DataPhaseInfo)
	s += fmt.Sprintf("OperationCode    : 0x%04x\n", o.OperationCode)
	s += fmt.Sprintf("TransactionID    : 0x%08x\n", o.TransactionID)
	s += fmt.Sprintf("P1               : 0x%08x\n", o.P1)
	s += fmt.Sprintf("P2               : 0x%08x\n", o.P2)
	s += fmt.Sprintf("P3               : 0x%08x\n", o.P3)
	s += fmt.Sprintf("P4               : 0x%08x\n", o.P4)
	return s
}

// OperationResponsePacket ...
type OperationResponsePacket struct {
	ResponseCode  uint16
	TransactionID uint32
	P1            uint32
	P2            uint32
	P3            uint32
	P4            uint32
}

func (o OperationResponsePacket) String() string {
	var s string
	s += fmt.Sprintf("----------------\n")
	s += fmt.Sprintf("OperationResponsePacket\n")
	s += fmt.Sprintf("----------------\n")
	s += fmt.Sprintf("ResponseCode     : 0x%08x\n", o.ResponseCode)
	s += fmt.Sprintf("TransactionID    : %v\n", o.TransactionID)
	s += fmt.Sprintf("Parameter1       : %v\n", o.P1)
	s += fmt.Sprintf("Parameter2       : %v\n", o.P2)
	s += fmt.Sprintf("Parameter3       : %v\n", o.P3)
	s += fmt.Sprintf("Parameter4       : %v\n", o.P4)
	return s
}

// EventPacket ...
type EventPacket struct {
	EventCode     uint16
	TransactionID uint32
	P1            uint32
	P2            uint32
	P3            uint32
}
