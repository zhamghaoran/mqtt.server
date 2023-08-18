package packets

import (
	"fmt"
	"io"
	"leetcode/constant"
)

type PubackPacket struct {
	FixedHeader
	MessageID uint16
}

func (pa *PubackPacket) Type() int {
	return constant.PUBACK
}

func (pa *PubackPacket) String() string {
	return fmt.Sprintf("%s MessageID: %d", pa.FixedHeader, pa.MessageID)
}

func (pa *PubackPacket) Write(w io.Writer) error {
	var err error
	pa.FixedHeader.RemainingLength = 2
	packet := pa.FixedHeader.pack()
	packet.Write(encodeUint16(pa.MessageID))
	_, err = packet.WriteTo(w)

	return err
}

func (pa *PubackPacket) Unpack(b io.Reader) error {
	var err error
	pa.MessageID, err = decodeUint16(b)

	return err
}

func (pa *PubackPacket) Details() Details {
	return Details{Qos: pa.Qos, MessageID: pa.MessageID}
}
