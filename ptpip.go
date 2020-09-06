package ptpip

import (
	"fmt"
	"net"

	"github.com/takurooo/ptpip/packet"
)

const (
	port string = ":15740"
)

// Initiator ...
type Initiator struct {
	GUID            []byte
	FriendlyName    string
	ProtocolVersion uint32
}

// Client ...
type Client struct {
	cConn net.Conn
	eConn net.Conn
	host  string
	ini   *Initiator
	stop  chan struct{}
	done  chan struct{}
}

func (c *Client) eventReciever() {
	defer func() {
		close(c.done)
	}()

	for {
		eventCode, err := packet.RecvEvent(c.eConn)
		fmt.Println(eventCode, err)

		select {
		case <-c.stop:
			return
		default:
		}

	}
}

// NewClient ...
func NewClient(host string, initiator *Initiator) *Client {
	if initiator == nil {
		initiator = new(Initiator)
		initiator.GUID = []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0A, 0x0B, 0x0C, 0xD, 0xE, 0x0F}
		initiator.FriendlyName = "hogehoge"
		initiator.ProtocolVersion = uint32(0x00010000)
	}
	return &Client{cConn: nil, eConn: nil, host: host, ini: initiator, stop: make(chan struct{}), done: make(chan struct{})}
}

// Disconnect ...
func (c *Client) Disconnect() (err error) {

	if err = c.cConn.Close(); err != nil {
		return err
	}
	if err = c.eConn.Close(); err != nil {
		return err
	}

	// TCPのコネクションを閉じないとgoroutineがTCPのリード待ちから返ってこれないので
	// TCPのコネクションを閉じてからchannelで終了指示を送る
	close(c.stop)
	<-c.done

	return nil
}

// Connect ...
func (c *Client) Connect() (err error) {
	addr := c.host + port
	// ---------------------------------------
	// establish connection for ptp-ip command
	// ---------------------------------------
	c.cConn, err = net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	initCommandRequestPacket := &(packet.InitCommandRequestPacket{
		GUID:            c.ini.GUID,
		FriendlyName:    c.ini.FriendlyName,
		ProtocolVersion: c.ini.ProtocolVersion,
	})

	ackPacket, err := packet.InitCommandRequest(c.cConn, initCommandRequestPacket)
	if err != nil {
		return err
	}

	// ---------------------------------------
	// establish connection for ptp-ip event
	// ---------------------------------------
	c.eConn, err = net.Dial("tcp", addr)
	if err != nil {
		return err
	}

	err = packet.InitEventRequest(c.eConn, ackPacket.ConnectionNumber)
	if err != nil {
		return err
	}

	go c.eventReciever()

	return nil
}

// OperationRequest ...
func (c *Client) OperationRequest(opCode uint16, phase uint32, transactionID uint32, p1, p2, p3, p4 uint32, sendData []byte) (recvData []byte, err error) {

	req := &packet.OperationRequestPacket{
		DataPhaseInfo: phase,
		OperationCode: opCode,
		TransactionID: transactionID,
		P1:            p1,
		P2:            p2,
		P3:            p3,
		P4:            p4,
	}

	recvData, err = packet.OperationRequest(c.cConn, req, sendData)
	if err != nil {
		return nil, err
	}

	return recvData, err
}
