package packets

import (
	"bytes"
	"fmt"
	"io"
	"mqtt-server/constant"
)

type ConnackPacket struct {
	FixedHeader
	SessionPresent bool
	ReturnCode     byte
}

func (ca *ConnackPacket) Type() int {
	return constant.CONNACK
}

func (ca *ConnackPacket) String() string {
	return fmt.Sprintf("%s sessionpresent: %t returncode: %d", ca.FixedHeader, ca.SessionPresent, ca.ReturnCode)
}

func (ca *ConnackPacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error

	body.WriteByte(boolToByte(ca.SessionPresent))
	body.WriteByte(ca.ReturnCode)
	ca.FixedHeader.RemainingLength = 2
	packet := ca.FixedHeader.pack()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(w)

	return err
}

func (ca *ConnackPacket) Unpack(b io.Reader) error {
	flags, err := decodeByte(b)
	if err != nil {
		return err
	}
	ca.SessionPresent = 1&flags > 0
	ca.ReturnCode, err = decodeByte(b)

	return err
}

func (ca *ConnackPacket) Details() Details {
	return Details{Qos: 0, MessageID: 0, Address: ca.RemoteAddress}
}
