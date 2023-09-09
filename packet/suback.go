package packets

import (
	"bytes"
	"fmt"
	"io"
	"mqtt/constant"
)

type SubackPacket struct {
	FixedHeader
	MessageID   uint16
	ReturnCodes []byte
}

func (sa *SubackPacket) Type() int {
	return constant.SUBACK
}

func (sa *SubackPacket) String() string {
	return fmt.Sprintf("%s MessageID: %d", sa.FixedHeader, sa.MessageID)
}

func (sa *SubackPacket) Write(w io.Writer) error {
	var body bytes.Buffer
	var err error
	body.Write(encodeUint16(sa.MessageID))
	body.Write(sa.ReturnCodes)
	sa.FixedHeader.RemainingLength = body.Len()
	packet := sa.FixedHeader.pack()
	packet.Write(body.Bytes())
	_, err = packet.WriteTo(w)

	return err
}

func (sa *SubackPacket) Unpack(b io.Reader) error {
	var qosBuffer bytes.Buffer
	var err error
	sa.MessageID, err = decodeUint16(b)
	if err != nil {
		return err
	}

	_, err = qosBuffer.ReadFrom(b)
	if err != nil {
		return err
	}
	sa.ReturnCodes = qosBuffer.Bytes()

	return nil
}

func (sa *SubackPacket) Details() Details {
	return Details{Qos: 0, MessageID: sa.MessageID, Address: sa.RemoteAddress}
}
