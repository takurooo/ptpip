package packet

import (
	"bytes"
	"errors"
	"fmt"
	"io"

	"github.com/takurooo/binaryio"
	"github.com/takurooo/ptpip/buff"
)

const (
	endian = binaryio.LittleEndian
)

const (
	packetHeaderSize uint32 = 8
)

func dump(b []byte, col int) {
	for i := 0; i < col; i++ {
		fmt.Printf("%4x ", i)
	}
	fmt.Printf("\n")
	for i := 0; i < len(b); i++ {
		if i != 0 && i%col == 0 {
			fmt.Printf("\n")
		}
		fmt.Printf("0x%02x ", b[i])
	}
	fmt.Printf("\n")
}

func encodeFriendlyName(s string) []byte {
	buf := make([]byte, len(s)*2+2)
	for i, r := range s {
		buf[i*2] = byte(r)
		buf[i*2+1] = 0
	}
	buf[len(s)*2] = 0
	buf[len(s)*2+1] = 0
	return buf
}

func decodeFriendlyName(br *binaryio.Reader) string {
	frinedlyName := make([]uint8, 40)
	var i int
	for {
		v := br.ReadU16(endian)
		if v == 0x00 { // null terminated
			break
		}
		frinedlyName[i] = byte(v & 0xff)
		i++
	}

	return string(frinedlyName)
}

func sendPakcet(w io.Writer, packet []byte) (err error) {

	// fmt.Println("----------------")
	// fmt.Println("sendPakcet")
	// fmt.Println("----------------")
	// dump(packet, 4)

	_, err = w.Write(packet)
	if err != nil {
		return err
	}
	return nil
}

func recvPaket(r io.Reader) (packetLen uint32, packetType uint32, packetBody []byte, err error) {
	// read packet header
	var packetHeader = make([]byte, packetHeaderSize)
	_, err = r.Read(packetHeader)
	if err != nil {
		return 0, 0, nil, err
	}
	brHeader := binaryio.NewReader(bytes.NewReader(packetHeader))
	packetLen = brHeader.ReadU32(endian)
	packetType = brHeader.ReadU32(endian)

	// read packet body
	packetBodyLen := packetLen - packetHeaderSize
	if 0 < packetBodyLen {
		packetBody = make([]byte, packetBodyLen)
		_, err = r.Read(packetBody)
		if err != nil {
			return 0, 0, nil, err
		}
	}

	// fmt.Println("----------------")
	// fmt.Println("recvPaket")
	// fmt.Println("----------------")
	// dump(append(packetHeader, packetBody...), 4)

	return packetLen, packetType, packetBody, nil
}

func sendInitCommandRequestPacket(w io.Writer, p *InitCommandRequestPacket) (err error) {

	// check value
	if 16 < len(p.GUID) {
		return fmt.Errorf("invalid initiator GUID len 16 < %d", len(p.GUID))
	}
	if 19 < len(p.FriendlyName) {
		return fmt.Errorf("invalid initiator FriendlyName len 18 < %d", len(p.FriendlyName))
	}

	encodedFriendlyName := encodeFriendlyName(p.FriendlyName)
	packetLen := uint32(12 + len(p.GUID) + len(encodedFriendlyName))

	buff := buff.NewBuffer(int(packetLen))
	bw := binaryio.NewWriter(buff)

	// write packet header to buffer
	bw.WriteU32(packetLen, endian)
	bw.WriteU32(PacketTypeInitCommandRequest, endian)
	// write packet body to buffer
	bw.WriteRaw(p.GUID)
	bw.WriteRaw(encodedFriendlyName)
	bw.WriteU32(p.ProtocolVersion, endian)

	if bw.Err() != nil {
		return bw.Err()
	}

	fmt.Println("----------------")
	fmt.Println("sendInitCommandRequest")
	fmt.Println("----------------")
	fmt.Printf("packetLen        : 0x%08x\n", packetLen)
	fmt.Println(p)

	packet := buff.Bytes()

	err = sendPakcet(w, packet)
	if err != nil {
		return err
	}

	return nil
}

func recvInitCommandAckPacket(r io.Reader) (ack *InitCommandAckPacket, err error) {

	// read packet header
	packetLen, packetType, packetBody, err := recvPaket(r)
	if err != nil {
		return nil, err
	}

	if packetType != PacketTypeInitCommandAck {
		return nil, fmt.Errorf("invalid packet type 0x%08x expected 0x%08x", packetType, PacketTypeInitCommandAck)
	}

	// parse InitCommandAckPacket
	brBody := binaryio.NewReader(bytes.NewReader(packetBody))

	ack = &InitCommandAckPacket{}
	ack.ConnectionNumber = brBody.ReadU32(endian)
	ack.GUID = brBody.ReadRaw(16)
	ack.FriendlyName = decodeFriendlyName(brBody)
	ack.ProtocolVersion = brBody.ReadU32(endian)

	fmt.Println("----------------")
	fmt.Println("recvInitCommandAck")
	fmt.Println("----------------")
	fmt.Printf("packetLen        : 0x%08x\n", packetLen)
	fmt.Printf("packetType       : 0x%08x\n", packetType)
	fmt.Println(ack)

	return ack, nil
}

func sendInitEventRequestPacket(w io.Writer, conndectionNumber uint32) (err error) {

	packetLen := uint32(12)

	buff := buff.NewBuffer(int(packetLen))
	bw := binaryio.NewWriter(buff)

	// write packet header to buffer
	bw.WriteU32(packetLen, endian)
	bw.WriteU32(PacketTypeInitEventRequest, endian)
	// write packet body to buffer
	bw.WriteU32(conndectionNumber, endian)

	if bw.Err() != nil {
		return bw.Err()
	}

	fmt.Println("----------------")
	fmt.Println("sendInitEventRequest")
	fmt.Println("----------------")
	fmt.Printf("packetLen        : 0x%08x\n", packetLen)
	fmt.Printf("ConnectionNumber : 0x%08x\n", conndectionNumber)

	packet := buff.Bytes()

	err = sendPakcet(w, packet)
	if err != nil {
		return err
	}

	return nil
}

func recvInitEventAckPacket(r io.Reader) error {

	// read packet header
	packetLen, packetType, _, err := recvPaket(r)
	if err != nil {
		return err
	}

	if packetType != PacketTypeInitEventAck {
		return fmt.Errorf("invalid packet type 0x%08x expected 0x%08x", packetType, PacketTypeInitEventAck)
	}

	fmt.Println("----------------")
	fmt.Println("recvInitEventAck")
	fmt.Println("----------------")
	fmt.Printf("packetLen        : 0x%08x\n", packetLen)
	fmt.Printf("packetType       : 0x%08x\n", packetType)

	return nil
}

func sendOperationRequestPacket(w io.Writer, req *OperationRequestPacket) (err error) {

	packetLen := uint32(34)
	buff := buff.NewBuffer(int(packetLen))
	bw := binaryio.NewWriter(buff)

	// write packet header to buffer
	bw.WriteU32(packetLen, endian)
	bw.WriteU32(PacketTypeOperationRequest, endian)
	// write packet body to buffer
	bw.WriteU32(req.DataPhaseInfo, endian)
	bw.WriteU16(req.OperationCode, endian)
	bw.WriteU32(req.TransactionID, endian)
	bw.WriteU32(req.P1, endian)
	bw.WriteU32(req.P2, endian)
	bw.WriteU32(req.P3, endian)
	bw.WriteU32(req.P4, endian)

	if bw.Err() != nil {
		return bw.Err()
	}
	fmt.Println("----------------")
	fmt.Println("sendOperationRequestPacket")
	fmt.Println("----------------")
	fmt.Printf("packetLen        : 0x%08x\n", packetLen)
	fmt.Println(req)

	packet := buff.Bytes()

	err = sendPakcet(w, packet)
	if err != nil {
		return err
	}

	return nil
}

func sendDataPacket(w io.Writer, sendData []byte) (err error) {

	if len(sendData) == 0 {
		return errors.New("send data empty")
	}

	{
		packetLen := uint32(20)
		buff := buff.NewBuffer(int(packetLen))
		bw := binaryio.NewWriter(buff)

		// write packet header to buffer
		bw.WriteU32(packetLen, endian)
		bw.WriteU32(PacketTypeStartData, endian)
		// write packet body to buffer
		bw.WriteU32(1, endian)
		bw.WriteU64(uint64(len(sendData)), endian)

		if bw.Err() != nil {
			return bw.Err()
		}

		packet := buff.Bytes()

		err = sendPakcet(w, packet)
		if err != nil {
			return err
		}
	}
	{
		packetLen := uint32(12 + len(sendData))
		buff := buff.NewBuffer(int(packetLen))
		bw := binaryio.NewWriter(buff)

		// write packet header to buffer
		bw.WriteU32(packetLen, endian)
		bw.WriteU32(PacketTypeData, endian)
		// write packet body to buffer
		bw.WriteU32(1, endian)
		bw.WriteRaw(sendData)

		if bw.Err() != nil {
			return bw.Err()
		}

		packet := buff.Bytes()

		err = sendPakcet(w, packet)
		if err != nil {
			return err
		}
	}
	{
		packetLen := uint32(12)
		buff := buff.NewBuffer(int(packetLen))
		bw := binaryio.NewWriter(buff)

		// write packet header to buffer
		bw.WriteU32(packetLen, endian)
		bw.WriteU32(PacketTypeEndData, endian)
		// write packet body to buffer
		bw.WriteU32(1, endian)

		if bw.Err() != nil {
			return bw.Err()
		}

		packet := buff.Bytes()

		err = sendPakcet(w, packet)
		if err != nil {
			return err
		}
	}

	return nil
}

func recvDataPacket(r io.Reader) (data []byte, err error) {
	var (
		packetLen       uint32
		packetType      uint32
		packetBody      []byte
		totalDataLength uint64
	)

	buff := buff.NewBuffer(64)

L:
	for {
		packetLen, packetType, packetBody, err = recvPaket(r)
		brBody := binaryio.NewReader(bytes.NewReader(packetBody))
		var payload []byte

		switch packetType {
		case PacketTypeStartData:
			// fmt.Println("PacketTypeStartData")
			_ = brBody.ReadU32(endian) // transactionID
			totalDataLength = brBody.ReadU64(endian)
		case PacketTypeData:
			// fmt.Println("PacketTypeData")
			_ = brBody.ReadU32(endian)                       // transactionID
			payload = brBody.ReadRaw(uint64(packetLen - 12)) // packetLen - packetHeaderSize - 4
			buff.Write(payload)
		case PacketTypeEndData:
			// fmt.Println("PacketTypeEndData")
			_ = brBody.ReadU32(endian)                       // transactionID
			payload = brBody.ReadRaw(uint64(packetLen - 12)) // packetLen - packetHeaderSize - 4
			buff.Write(payload)
			break L
		}
	}

	if uint64(buff.Len()) != totalDataLength {
		return nil, fmt.Errorf("invalid data len 0x%x expected 0x%x", buff.Len(), totalDataLength)
	}

	return buff.Bytes(), nil
}

func recvOperationReponsePacket(r io.Reader) (resp *OperationResponsePacket, err error) {

	// read packet header
	packetLen, packetType, packetBody, err := recvPaket(r)
	if err != nil {
		return nil, err
	}

	if packetType != PacketTypeOperationResponse {
		return nil, fmt.Errorf("invalid packet type 0x%08x expected 0x%08x", packetType, PacketTypeOperationResponse)
	}

	// parse OperationResponsePacket
	brBody := binaryio.NewReader(bytes.NewReader(packetBody))

	resp = &OperationResponsePacket{}
	resp.ResponseCode = brBody.ReadU16(endian)
	resp.TransactionID = brBody.ReadU32(endian)
	resp.P1 = brBody.ReadU32(endian)
	resp.P2 = brBody.ReadU32(endian)
	resp.P3 = brBody.ReadU32(endian)
	resp.P4 = brBody.ReadU32(endian)

	fmt.Println("----------------")
	fmt.Println("recvOperationReponsePacket")
	fmt.Println("----------------")
	fmt.Printf("packetLen        : 0x%08x\n", packetLen)
	fmt.Printf("packetType       : 0x%08x\n", packetType)
	fmt.Println(resp)

	return resp, nil
}

func parseEventPacket(packetBody []byte) (e *EventPacket, err error) {

	// parse EventPacket
	brBody := binaryio.NewReader(bytes.NewReader(packetBody))

	e = &EventPacket{}
	e.EventCode = brBody.ReadU16(endian)
	e.TransactionID = brBody.ReadU32(endian)
	e.P1 = brBody.ReadU32(endian)
	e.P2 = brBody.ReadU32(endian)
	e.P3 = brBody.ReadU32(endian)

	return e, nil
}

func sendProbeResponsePacket(w io.Writer) (err error) {

	packetLen := uint32(8)
	buff := buff.NewBuffer(int(packetLen))
	bw := binaryio.NewWriter(buff)

	// write packet header to buffer
	bw.WriteU32(packetLen, endian)
	bw.WriteU32(PacketTypeProbeResponse, endian)

	if bw.Err() != nil {
		return bw.Err()
	}

	packet := buff.Bytes()

	err = sendPakcet(w, packet)
	if err != nil {
		return err
	}

	return nil
}

// InitCommandRequest ...
func InitCommandRequest(conn PTPIPConn, p *InitCommandRequestPacket) (ack *InitCommandAckPacket, err error) {

	err = sendInitCommandRequestPacket(conn, p)
	if err != nil {
		return nil, err
	}
	ack, err = recvInitCommandAckPacket(conn)
	if err != nil {
		return nil, err
	}

	return ack, nil
}

// InitEventRequest ...
func InitEventRequest(conn PTPIPConn, conndectionNumber uint32) (err error) {
	err = sendInitEventRequestPacket(conn, conndectionNumber)
	if err != nil {
		return err
	}
	err = recvInitEventAckPacket(conn)
	if err != nil {
		return err
	}
	return nil
}

// OperationRequest ...
func OperationRequest(conn PTPIPConn, req *OperationRequestPacket, sendData []byte) (recvData []byte, err error) {

	err = sendOperationRequestPacket(conn, req)
	if err != nil {
		return nil, err
	}

	switch req.DataPhaseInfo {
	case DataPhaseInfoNoDataOrDataIn:
		recvData, err = recvDataPacket(conn)
		if err != nil {
			return nil, err
		}
	case DataPhaseInfoDataOut:
		err = sendDataPacket(conn, sendData)
		if err != nil {
			return nil, err
		}
	}

	resp, err := recvOperationReponsePacket(conn)
	if err != nil {
		return nil, err
	}

	if resp.ResponseCode != ResponseCodeOK {
		return nil, fmt.Errorf("operation response error 0x%08x", resp.ResponseCode)
	}

	return recvData, nil
}

// RecvEvent ...
func RecvEvent(conn PTPIPConn) (eventCode uint16, err error) {

	var e *EventPacket
L:
	for {
		// read packet header
		_, packetType, packetBody, err := recvPaket(conn)
		if err != nil {
			return 0, err
		}

		switch packetType {
		case PacketTypeEvent:
			e, err = parseEventPacket(packetBody)
			if err != nil {
				return 0, err
			}
			break L
		case PacketTypeProbeRequest:
			err = sendProbeResponsePacket(conn)
			if err != nil {
				return 0, err
			}
		}
	}

	return e.EventCode, nil
}
